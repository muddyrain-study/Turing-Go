package controller

import (
	"Turing-Go/config"
	"Turing-Go/constant"
	"Turing-Go/net"
	"log"
	"strings"
	"sync"
)

var GateHandler = &Handler{
	proxyMap: make(map[string]map[int64]*net.ProxyClient),
}

type Handler struct {
	proxyMutex sync.Mutex
	// 代理地址 => 客户端连接 => 代理客户端
	proxyMap   map[string]map[int64]*net.ProxyClient
	loginProxy string
	gameProxy  string
}

func (h *Handler) Router(r *net.Router) {
	h.loginProxy = config.File.MustValue("gate_server", "login_proxy", "ws://127.0.0.1:8003")
	h.gameProxy = config.File.MustValue("gate_server", "game_proxy", "ws://127.0.0.1:8081")
	g := r.Group("*")
	g.AddRouter("*", h.all)
}

func (h *Handler) all(req *net.WsMsgReq, resp *net.WsMsgResp) {
	log.Println("req body name", req.Body.Name)
	proxyStr := req.Body.Proxy
	if isAccount(req.Body.Name) {
		proxyStr = h.loginProxy
	} else {
		proxyStr = h.gameProxy
	}
	if proxyStr == "" {
		resp.Body.Code = constant.ProxyNotInConnect
		return
	}
	h.proxyMutex.Lock()
	_, ok := h.proxyMap[proxyStr]
	if !ok {
		h.proxyMap[proxyStr] = make(map[int64]*net.ProxyClient)
	}
	h.proxyMutex.Unlock()
	//获取客户端id
	cidValue, _ := req.Conn.GetProperty("cid")
	cid := cidValue.(int64)
	proxy, ok := h.proxyMap[proxyStr][cid]
	if !ok {
		//没有 建立连接 并发起调用
		proxy = net.NewProxyClient(proxyStr)
		h.proxyMutex.Lock()
		h.proxyMap[proxyStr][cid] = proxy
		h.proxyMutex.Unlock()
		err := proxy.Connect()
		if err != nil {
			h.proxyMutex.Lock()
			delete(h.proxyMap[proxyStr], cid)
			h.proxyMutex.Unlock()
			resp.Body.Code = constant.ProxyConnectError
			log.Println("代理连接失败", err)
			return
		}
		proxy.SetProperty("cid", cid)
		proxy.SetProperty("proxy", proxyStr)
		proxy.SetProperty("gateConn", req.Conn)
		proxy.SetOnPush(h.onPush)
	}
	resp.Body.Seq = req.Body.Seq
	resp.Body.Name = req.Body.Name
	r, err := proxy.Send(req.Body.Name, req.Body.Msg)
	if err == nil {
		resp.Body.Code = r.Code
		resp.Body.Msg = r.Msg
	} else {
		resp.Body.Code = constant.ProxyConnectError
		resp.Body.Msg = nil
	}
}

func (h *Handler) onPush(conn *net.ClientConn, body *net.RespBody) {
	gateConn, err := conn.GetProperty("gateConn")
	if err != nil {
		return
	}
	wc := gateConn.(net.WSConn)
	wc.Push(body.Name, body.Msg)
}

func isAccount(name string) bool {
	return strings.HasPrefix(name, "account.")
}
