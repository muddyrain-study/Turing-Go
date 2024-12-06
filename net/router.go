package net

import (
	"log"
	"strings"
	"sync"
)

type Handler func(req *WsMsgReq, resp *WsMsgResp)
type MiddlewareFunc func(handler Handler) Handler
type group struct {
	mutex             sync.RWMutex
	prefix            string
	handlerMap        map[string]Handler
	MiddlewareFuncMap map[string][]MiddlewareFunc
	Middlewares       []MiddlewareFunc
}

func (g *group) exec(name string, req *WsMsgReq, resp *WsMsgResp) {
	h, ok := g.handlerMap[name]
	if !ok {
		h, ok = g.handlerMap["*"]
		if !ok {
			log.Println("路由未定义")
		}
		if h != nil {
			h(req, resp)
		}
	}
	else{
		for i := 0; i < len(g.Middlewares); i++ {
			h  = g.Middlewares[i](h)
		}
		mm,ok := g.MiddlewareFuncMap[name]
		if ok {
			for i := 0; i < len(g.Middlewares); i++ {
						h = mm[i](h)
			}
		}
		if h != nil {
			h(req, resp)
		}
	}

}
func (g *group) AddRouter(name string, h Handler, middlewares ...MiddlewareFunc) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.handlerMap[name] = h
	g.MiddlewareFuncMap[name] = middlewares
}
func (g *group) Use(middlewares ...MiddlewareFunc) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.Middlewares = append(g.Middlewares, middlewares...)
}
func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix:     prefix,
		handlerMap: map[string]Handler{},
		MiddlewareFuncMap: make(map[string][]MiddlewareFunc),
	}

	r.group = append(r.group, g)
	return g
}

type Router struct {
	group []*group
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Run(req *WsMsgReq, resp *WsMsgResp) {
	strs := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""
	if len(strs) == 2 {
		prefix = strs[0]
		name = strs[1]
	}
	for _, g := range r.group {
		if g.prefix == prefix {
			g.exec(name, req, resp)
		} else if g.prefix == "*" {
			g.exec(name, req, resp)
		}
	}
}
