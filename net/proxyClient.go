package net

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type ProxyClient struct {
	proxy string
	conn  *ClientConn
}

func (c *ProxyClient) Connect() error {
	// 连接 ws 服务器
	//通过Dialer连接websocket服务器
	var dialer = websocket.Dialer{
		Subprotocols:     []string{"p1", "p2"},
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 30 * time.Second,
	}
	ws, _, err := dialer.Dial(c.proxy, nil)
	if err == nil {
		c.conn = NewClientConn(ws)
		if !c.conn.Start() {
			return errors.New("握手失败")
		}
		log.Println("连接成功 proxy:", c.proxy)
	}
	return err
}

func (c *ProxyClient) SetProperty(s string, data interface{}) {
	if c.conn != nil {
		c.conn.SetProperty(s, data)
	}
}

func (c *ProxyClient) SetOnPush(hook func(conn *ClientConn, body *RespBody)) {
	if c.conn != nil {
		c.conn.SetOnPush(hook)
	}
}

func (c *ProxyClient) Send(name string, msg interface{}) (*RespBody, error) {
	if c.conn != nil {
		return c.conn.Send(name, msg)
	}
	return nil, errors.New("连接未发现")
}

func NewProxyClient(proxy string) *ProxyClient {
	return &ProxyClient{
		proxy: proxy,
	}
}
