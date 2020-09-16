package main

import (
	"fmt"
	"github.com/dawnbyte/go-sdk/client"
	"net/http"
	_ "net/http/pprof"

)


func main() {

	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()

	//初始化Mega sdk
	client.New().
		WithId(1).
		WithWsAddress("wss://stream.dawnbyte.com/ws").
		RegisterWsMessageHandler(&handler{}).
		Run()

}

//平台消息接收类
//实现接收方法 传入接收类 即可接收到平台回调
type handler struct{}

//消息接收方法
func (h *handler) ProcessMsg(data []byte) {
	//这里接收到ws数据 商户实现自己的业务逻辑
	fmt.Println(string(data))
}
