package controller

import (
	"Turing-Go/constant"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/logic"
	"Turing-Go/server/game/middleware"
	"Turing-Go/server/game/model"
)

var UnionController = &unionController{}

type unionController struct{}

func (c *unionController) InitRouter(router *net.Router) {
	g := router.Group("union")
	g.Use(middleware.Log())
	g.AddRouter("list", c.list, middleware.CheckRole())
}

func (c *unionController) list(req *net.WsMsgReq, resp *net.WsMsgResp) {
	respObj := &model.ListRsp{}
	// 查询数据库查询所有联盟
	resp.Body.Code = constant.OK
	resp.Body.Msg = respObj

	uns, err := logic.CoalitionService.List()
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	respObj.List = uns
}
