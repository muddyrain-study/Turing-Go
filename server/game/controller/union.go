package controller

import (
	"Turing-Go/constant"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/logic"
	"Turing-Go/server/game/middleware"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"github.com/mitchellh/mapstructure"
)

var UnionController = &unionController{}

type unionController struct{}

func (c *unionController) InitRouter(router *net.Router) {
	g := router.Group("union")
	g.Use(middleware.Log())
	g.AddRouter("list", c.list, middleware.CheckRole())
	g.AddRouter("info", c.info, middleware.CheckRole())
	g.AddRouter("applyList", c.applyList, middleware.CheckRole())
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

func (c *unionController) info(req *net.WsMsgReq, resp *net.WsMsgResp) {
	reqObj := &model.InfoReq{}
	err := mapstructure.Decode(req.Body.Msg, reqObj)
	if err != nil {
		resp.Body.Code = constant.InvalidParam
		return
	}
	respObj := &model.InfoRsp{}
	resp.Body.Msg = respObj
	resp.Body.Code = constant.OK
	un, err := logic.CoalitionService.Get(reqObj.Id)
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	respObj.Info = un
	respObj.Id = reqObj.Id
}

func (c *unionController) applyList(req *net.WsMsgReq, resp *net.WsMsgResp) {
	//根据联盟id 去查询申请列表，rid申请人，你角色表 查询详情即可
	// state 0 正在申请 1 拒绝 2 同意
	//什么人能看到申请列表 只有盟主和副盟主能看到申请列表

	reqObj := &model.ApplyReq{}
	err := mapstructure.Decode(req.Body.Msg, reqObj)
	if err != nil {
		resp.Body.Code = constant.InvalidParam
		return
	}
	rspObj := &model.ApplyRsp{}
	resp.Body.Code = constant.OK
	resp.Body.Msg = rspObj

	r, _ := req.Conn.GetProperty("role")
	role := r.(data.Role)

	//查询联盟
	un := logic.CoalitionService.GetCoalition(reqObj.Id)
	if un == nil {
		resp.Body.Code = constant.DBError
		return
	}

	if un.Chairman != role.RId && un.ViceChairman != role.RId {
		rspObj.Id = reqObj.Id
		rspObj.Applys = make([]model.ApplyItem, 0)
		return
	}
	ais, err := logic.CoalitionService.GetListApply(reqObj.Id, 0)
	if err != nil {
		resp.Body.Code = constant.DBError
		return
	}
	rspObj.Id = reqObj.Id
	rspObj.Applys = ais
}
