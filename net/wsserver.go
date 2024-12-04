package net

import (
	"Turing-Go/utils"
	"encoding/json"
	"errors"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"sync"
	"time"
)

var cid int64

type wsServer struct {
	wsConn       *websocket.Conn
	router       *Router
	outChan      chan *WsMsgResp
	seq          int
	property     map[string]interface{}
	propertyLock sync.RWMutex
	needSecret   bool
}

func NewWsServer(wsConn *websocket.Conn, needSecret bool) *wsServer {
	s := &wsServer{
		wsConn:     wsConn,
		outChan:    make(chan *WsMsgResp, 1000),
		property:   make(map[string]interface{}),
		seq:        0,
		needSecret: needSecret,
	}
	cid++
	s.SetProperty("cid", cid)
	return s
}

func (w *wsServer) Router(router *Router) {
	w.router = router
}

func (w *wsServer) SetProperty(key string, value interface{}) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}
func (w *wsServer) GetProperty(key string) (interface{}, error) {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()
	if v, ok := w.property[key]; ok {
		return v, nil
	} else {
		return nil, errors.New("no property found")
	}
}
func (w *wsServer) RemoveProperty(key string) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	delete(w.property, key)
}
func (w *wsServer) Addr() string {
	return w.wsConn.RemoteAddr().String()
}
func (w *wsServer) Push(name string, data interface{}) {
	resp := &WsMsgResp{
		Body: &RespBody{
			Name: name,
			Msg:  data,
			Seq:  0,
		},
	}
	w.outChan <- resp

}

func (w *wsServer) Start() {
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

func (w *wsServer) writeMsgLoop() {
	for {
		select {
		case msg := <-w.outChan:
			jsonData, _ := json.Marshal(msg.Body)
			log.Printf("服务器发送消息: %v \n", string(jsonData))
			w.Write(msg.Body)
		}
	}
}

func (w *wsServer) Write(msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	secretKey, err := w.GetProperty("secretKey")
	if err == nil {
		key := secretKey.(string)
		data, _ = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	}
	if data, err := utils.Zip(data); err == nil {
		if err := w.wsConn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			w.Close()
		}
	}
}

func (w *wsServer) readMsgLoop() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("readMsgLoop panic: ", r)
			w.Close()
		}
	}()
	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			log.Println("read message error: ", err)
			break
		}
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("Failed to unzip data: ", err)
			continue
		}
		if w.needSecret {
			secretKey, err := w.GetProperty("secretKey")
			if err == nil {
				key := secretKey.(string)
				d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
				if err != nil {
					log.Println("Failed to decrypt data: ", err)
					//continue
					w.Handshake()
				} else {
					data = d
				}
			}
		}
		body := &ReqBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("Failed to json unmarshal data: ", err)
			continue
		} else {
			req := &WsMsgReq{Conn: w, Body: body}
			resp := &WsMsgResp{Body: &RespBody{Name: body.Name, Seq: req.Body.Seq}}
			if req.Body.Name == "heartbeat" {
				// 回复心跳
				h := &Heartbeat{}
				err := mapstructure.Decode(resp.Body.Msg, h)
				if err != nil {
					log.Println("Failed to decode heartbeat: ", err)
				}
				h.STime = time.Now().UnixNano() / 1e6
				resp.Body.Msg = h
			} else {
				if w.router != nil {
					w.router.Run(req, resp)
					w.outChan <- resp
				}
			}

		}

	}
	w.Close()
}

func (w *wsServer) Close() {
	err := w.wsConn.Close()
	if err != nil {
		log.Println("Failed to close websocket connection: ", err)
	}
}

const handshakeMsg = "handshake"

func (w *wsServer) Handshake() {
	secretKey := ""
	key, err := w.GetProperty("secretKey")
	if err == nil {
		secretKey = key.(string)
	} else {
		secretKey = utils.RandSeq(16)
	}
	handshake := &Handshake{Key: secretKey}
	body := &RespBody{
		Name: handshakeMsg,
		Msg:  handshake,
	}

	if data, err := json.Marshal(body); err == nil {
		if secretKey != "" {
			w.SetProperty("secretKey", secretKey)
		} else {
			w.RemoveProperty("secretKey")
		}
		if data, err := utils.Zip(data); err == nil {
			w.wsConn.WriteMessage(websocket.BinaryMessage, data)
		}
	}
}
