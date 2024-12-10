package controller

import (
	"Turing-Go/net"
	"Turing-Go/server/game/middleware"
)

var CityController = &cityController{}

type cityController struct {
}

func (g *cityController) InitRouter(router *net.Router) {
	r := router.Group("city")
	r.Use(middleware.Log())
	r.AddRouter("facilities", g.facilities, middleware.CheckRole())
}

func (g *cityController) facilities(req *net.WsMsgReq, resp *net.WsMsgResp) {

}
