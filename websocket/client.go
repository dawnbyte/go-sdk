package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

type MegaResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type WsMessageHandler interface {
	ProcessMsg(data []byte)
}

//启动ws
func Run(url string, handler WsMessageHandler) {
	//读写协程失败后发送通知
	errChan := make(chan struct{})
	//主协程通知读写通知关闭 回收资源
	stop := make(chan struct{})
	//确保读写协程失败后只进行一次重启
	once := &sync.Once{}
	//等待读写资源回收后重启
	wg := &sync.WaitGroup{}

	go connect(url,handler,errChan,stop,once,wg)
	<- errChan
	close(stop)

	wg.Wait()
	log.Println("连接异常 开始重连...")
	time.Sleep(time.Millisecond*50)
	Run(url, handler)
}

func connect(url string,handler WsMessageHandler, errChan chan struct{}, stop chan struct{}, once *sync.Once, wg *sync.WaitGroup) {
	dial := websocket.Dialer{
		HandshakeTimeout: time.Second * 3,
	}
	//连接ws
	conn, response, err := dial.Dial(url, nil)
	if err != nil {
		//网络原因没有响应体(超时等) 此时可直接结束
		if response != nil {
			//有响应体说明服务端已处理连接,但是拒绝连接. 解析响应body来获取具体错误原因
			b, _ := ioutil.ReadAll(response.Body)
			//body结构反解析
			megaResponse := &MegaResponse{}
			_ = json.Unmarshal(b, megaResponse)

			//code除了0其他为异常情况  如  10000 Illegal IP
			if megaResponse.Code != 0 {
				log.Println(fmt.Sprintf("连接出错: %s", megaResponse.Message))
			}
		}
		log.Println("连接出错" + err.Error())
		errChan <- struct{}{}
		return
	}
	log.Println("连接成功!")
	wg.Add(2)
	//没有报错说明已经连接成功 此时得到可用ws链接
	go readPump(conn, handler, errChan, stop, wg, once)
	go writePump(conn, handler, errChan, stop, wg, once)
}

func readPump(conn *websocket.Conn, handler WsMessageHandler, errChan chan struct{}, stopChan chan struct{}, wg *sync.WaitGroup, once *sync.Once) {
	defer func() {
		_ = conn.Close()
		wg.Done()
	}()

	for {
		select {
		case _,ok := <- stopChan:
			if !ok {
				fmt.Println("readPump收到结束信号")
				return
			}
		default:
			//接收ws数据
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("read ws message error: " + err.Error())
				once.Do(func() {
					errChan <- struct{}{}
				})
				return
			}
			//回调代理 处理消息
			handler.ProcessMsg(msg)
		}
	}
}

func writePump(conn *websocket.Conn, handler WsMessageHandler, errChan chan struct{}, stopChan chan struct{}, wg *sync.WaitGroup, once *sync.Once) {

	ticker := time.NewTicker(time.Second * 5)
	defer func() {
		ticker.Stop()
		_ = conn.Close()
		wg.Done()
	}()

	for {
		select {
		case <-ticker.C:
			//50s写入一次心跳消息
			if err := conn.WriteMessage(websocket.TextMessage, []byte("Ping")); err != nil {
				log.Println("send message error: " + err.Error())
				once.Do(func() {
					errChan <- struct{}{}
				})
				return
			}
		case _,ok := <- stopChan:
			if !ok {
				fmt.Println("writePump收到结束信号")
				return
			}

		}
	}
}
