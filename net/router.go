package net

import (
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
	strings := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""
	if len(strings) == 2 {
		prefix = strings[0]
		name = strings[1]
	}
	for _, g := range r.group {
		if g.prefix == prefix {
			g.exec(name, req, resp)
		} else if g.prefix == "*" {
			g.exec(name, req, resp)
		}
	}
}
