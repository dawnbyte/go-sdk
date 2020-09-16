# Dawnbyte 接入SDK - Golang

## 实现功能
- websocket快速接入
- websocket断线重连机制

## 使用示例
```go
func main()  {

	//初始化Mega sdk
	client.New().
		WithId(1).
		WithWsAddress("wss://stream.dawnbyte.com/ws").
		RegisterWsMessageHandler(&handler{}).
		Run()

	//集成SDK之后 应该由商户自己的主程序控制
	select {}
} 

//平台消息接收类
//实现接收方法 传入接收类 即可接收到平台回调
type handler struct {}

//消息接收方法
func (h *handler) ProcessMsg(data []byte) {
	//这里接收到ws数据 商户实现自己的业务逻辑
	fmt.Println(string(data))
}
```

## Todo
- 基础数据初始化
- ws断线期间 通过推拉结合进行数据的完整性拉取策略
- 更多优化...