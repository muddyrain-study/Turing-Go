package controller

import (
	"Turing-Go/constant"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/logic"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"Turing-Go/utils"
	"github.com/mitchellh/mapstructure"
	"log"
)

var RoleController = &roleController{}

type roleController struct {
}

func (r *roleController) InitRouter(router *net.Router) {
	g := router.Group("role")
	g.AddRouter("enterServer", r.enterServer)
	g.AddRouter("myProperty", r.myProperty)
}

func (r *roleController) enterServer(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	reqObj := &model.EnterServerReq{}
	rspObj := &model.EnterServerRsp{}

	err := mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	if err != nil {
		log.Println("enterServer mapstructure.Decode error:", err)
		rsp.Body.Code = constant.InvalidParam
		return
	}
	session := reqObj.Session
	_, claims, err := utils.ParseToken(session)
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	uid := claims.Uid
	err = logic.RoleService.EnterServer(uid, rspObj, req.Conn)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}

func (r *roleController) myProperty(req *net.WsMsgReq, resp *net.WsMsgResp) {
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		resp.Body.Code = constant.SessionInvalid
		return
	}
	resp.Body.Seq = req.Body.Seq
	resp.Body.Name = req.Body.Name
	rid := role.(data.Role).RId
	respObj := &model.MyRolePropertyRsp{}
	// 查询资源
	respObj.RoleRes, err = logic.RoleService.GetRoleRes(rid)
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	// 城池
	respObj.Citys, err = logic.RoleCityService.GetRoleCities(rid)
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	// 建筑
	respObj.MRBuilds, err = logic.RoleBuildService.GetBuilds(rid)
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	// 军队
	respObj.Armys, err = logic.ArmyService.GetArmies(rid)
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	// 武将
	respObj.Generals, err = logic.GeneralService.GetGenerals(rid)
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}

	resp.Body.Code = constant.OK
	resp.Body.Msg = respObj
}
