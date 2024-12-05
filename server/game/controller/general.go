package controller

import (
	"Turing-Go/constant"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/logic"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
)

var GeneralController = &generalController{}

type generalController struct {
}

func (g *generalController) InitRouter(router *net.Router) {
	r := router.Group("general")
	r.AddRouter("myGenerals", g.myGenerals)
}

func (g *generalController) myGenerals(req *net.WsMsgReq, resp *net.WsMsgResp) {

	role, err := req.Conn.GetProperty("role")
	if err != nil {
		resp.Body.Code = constant.SessionInvalid
		return
	}
	respObj := &model.MyGeneralRsp{}
	resp.Body.Msg = respObj
	resp.Body.Code = constant.OK
	resp.Body.Seq = req.Body.Seq
	resp.Body.Name = req.Body.Name

	rid := role.(data.Role).RId

	generals, err := logic.GeneralService.GetGenerals(rid)

	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	respObj.Generals = generals
	resp.Body.Msg = respObj
}
