package net

import (
	"Turing-Go/constant"
	"Turing-Go/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"sync"
)

type ClientConn struct {
	wsConn        *websocket.Conn
	isClosed      bool
	property      map[string]interface{}
	propertyLock  sync.RWMutex
	Seq           int64
	handshake     bool
	handshakeChan chan bool
	onPush        func(conn *ClientConn, body *RespBody)
	onClose       func(conn *ClientConn)
	syncCtxMap    map[int64]*syncCtx
	syncCtxLock   sync.RWMutex
}

type syncCtx struct {
	//Goroutine 的上下文，包含 Goroutine 的运行状态、环境、现场等信息
	ctx     context.Context
	cancel  context.CancelFunc
	outChan chan *RespBody
}

func NewSyncCtx() *syncCtx {
	ctx, cancel := context.WithCancel(context.Background())
	return &syncCtx{
		ctx:     ctx,
		cancel:  cancel,
		outChan: make(chan *RespBody),
	}
}

func (s *syncCtx) wait() *RespBody {
	select {
	case <-s.ctx.Done():
		log.Println("proxy service response message timeout")
		return nil
	case body := <-s.outChan:
		return body
	}
}

func (c *ClientConn) Start() bool {
	c.handshake = false
	go c.wsReadLoop()
	return c.waitHandshake()
}

func (c *ClientConn) waitHandshake() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()
	select {
	case <-c.handshakeChan:
		log.Println("握手成功")
		return true
	case <-ctx.Done():
		log.Println("握手超时")
		return false
	}

}

func (c *ClientConn) wsReadLoop() {
	//for {
	//	_, data, err := c.wsConn.ReadMessage()
	//	log.Println(data)
	//	log.Println(err)
	//	c.handshake = true
	//	c.handshakeChan <- true
	//
	//}
	defer func() {
		if r := recover(); r != nil {
			log.Println("read loop panic: ", r)
			c.Close()
		}
	}()
	for {
		_, data, err := c.wsConn.ReadMessage()
		if err != nil {
			log.Println("read message error: ", err)
			break
		}
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("Failed to unzip data: ", err)
			continue
		}
		secretKey, err := c.GetProperty("secretKey")
		if err == nil {
			key := secretKey.(string)
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println("Failed to decrypt data: ", err)
			} else {
				data = d
			}
		}
		body := &RespBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("Failed to json unmarshal data: ", err)
			continue
		} else {
			if body.Seq == 0 {
				if body.Name == handshakeMsg {
					hs := &Handshake{}
					err := mapstructure.Decode(body.Msg, hs)
					if err != nil {
						log.Println(err)
					}
					if hs.Key != "" {
						c.SetProperty("secretKey", hs.Key)
					} else {
						c.RemoveProperty("secretKey")
					}
					c.handshake = true
					c.handshakeChan <- true
				} else {
					if c.onPush != nil {
						c.onPush(c, body)
					}
				}
			} else {
				c.syncCtxLock.RLock()
				if ctx, ok := c.syncCtxMap[body.Seq]; ok {
					ctx.outChan <- body
				} else {
					log.Println("no seq found")
				}
				c.syncCtxLock.RUnlock()
			}
		}

	}
	c.Close()
}

func (c *ClientConn) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}
func (c *ClientConn) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if v, ok := c.property[key]; ok {
		return v, nil
	} else {
		return nil, errors.New("no property found")
	}
	return nil, nil
}
func (c *ClientConn) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}

func (c *ClientConn) Close() {
	c.wsConn.Close()
}

func (c *ClientConn) Addr() string {
	return c.wsConn.RemoteAddr().String()
}
func (c *ClientConn) Push(name string, data interface{}) {
	resp := &WsMsgResp{
		Body: &RespBody{
			Name: name,
			Msg:  data,
			Seq:  0,
		},
	}
	//c.outChan <- resp
	c.write(resp.Body)

}

func (c *ClientConn) write(body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		return err
	}
	secretKey, err := c.GetProperty("secretKey")
	if err == nil {
		key := secretKey.(string)
		data, err = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
		if err != nil {
			log.Println("Failed to encrypt data: ", err)
			return err
		}
	}
	if data, err := utils.Zip(data); err == nil {
		if err := c.wsConn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			log.Println("Failed to write message: ", err)
			return err
		}
	} else {
		log.Println("Failed to zip data: ", err)
		return err
	}
	return nil
}

func (c *ClientConn) SetOnPush(hook func(conn *ClientConn, body *RespBody)) {
	c.onPush = hook
}

func (c *ClientConn) Send(name string, msg interface{}) (*RespBody, error) {
	c.Seq += 1
	seq := c.Seq
	sc := NewSyncCtx()
	c.syncCtxLock.Lock()
	c.syncCtxMap[seq] = sc
	c.syncCtxLock.Unlock()

	req := &ReqBody{Name: name, Msg: msg, Seq: seq}
	err := c.write(req)
	resp := &RespBody{Name: name, Seq: seq, Code: constant.OK}
	if err != nil {
		sc.cancel()
	} else {
		r := sc.wait()
		if r == nil {
			resp.Code = constant.ProxyConnectError
		} else {
			resp = r
		}
	}
	c.syncCtxLock.Lock()
	delete(c.syncCtxMap, seq)
	c.syncCtxLock.Unlock()
	return resp, err
}

func NewClientConn(wsConn *websocket.Conn) *ClientConn {
	return &ClientConn{
		wsConn: wsConn,
	}
}
