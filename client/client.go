package client

import (
	"github.com/dawnbyte/go-sdk/websocket"
)

type Client struct {
	token     string
	id        int64
	wsAddress string
	handler   websocket.WsMessageHandler
}

func New() *Client {
	return &Client{}
}

func (c *Client) WithId(id int64) *Client {
	c.id = id
	return c
}

func (c *Client) WithToken(token string) *Client {
	c.token = token
	return c
}

func (c *Client) WithWsAddress(address string) *Client {
	c.wsAddress = address
	return c
}

func (c *Client) RegisterWsMessageHandler(handler websocket.WsMessageHandler) *Client {
	c.handler = handler
	return c
}

func (c *Client) Run() {

	if c.handler == nil {
		panic("缺少数据处理handler")
	}

	if c.wsAddress == "" {
		panic("缺少ws地址")
	}

	//启动ws服务
	websocket.Run(c.wsAddress, c.handler)
}
