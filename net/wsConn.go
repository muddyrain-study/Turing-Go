package net

import "sync"

type ReqBody struct {
	Seq   int64       `json:"seq"`
	Name  string      `json:"name"`
	Msg   interface{} `json:"msg"`
	Proxy string      `json:"proxy"`
}

type RespBody struct {
	Seq  int64       `json:"seq"`
	Name string      `json:"name"`
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

type WsContext struct {
	mutex    sync.RWMutex
	property map[string]interface{}
}

func (w *WsContext) Set(name string, value interface{}) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.property[name] = value
}

func (w *WsContext) Get(name string) interface{} {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	value, ok := w.property[name]
	if !ok {
		return nil
	}
	return value

}

type WsMsgReq struct {
	Body    *ReqBody
	Conn    WSConn
	Context *WsContext
}

type WsMsgResp struct {
	Body *RespBody
}

type WSConn interface {
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
	Addr() string
	Push(name string, data interface{})
}

type Handshake struct {
	Key string `json:"key"`
}

type Heartbeat struct {
	CTime int64 `json:"ctime"`
	STime int64 `json:"stime"`
}
