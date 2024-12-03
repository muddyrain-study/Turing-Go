package net

import (
	"log"
	"strings"
)

type Handler func(req *WsMsgReq, resp *WsMsgResp)

type group struct {
	prefix     string
	handlerMap map[string]Handler
}

func (g *group) exec(name string, req *WsMsgReq, resp *WsMsgResp) {
	h := g.handlerMap[name]
	if h != nil {
		h(req, resp)
	} else {
		h = g.handlerMap["*"]
		if h != nil {
			h(req, resp)
		} else {
			log.Println("路由未定义")
		}
	}
}
func (g *group) AddRouter(name string, h Handler) {
	g.handlerMap[name] = h
}
func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix:     prefix,
		handlerMap: map[string]Handler{},
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
