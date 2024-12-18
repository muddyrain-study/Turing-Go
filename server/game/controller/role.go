package controller

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/logic"
	"Turing-Go/server/game/middleware"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"Turing-Go/utils"
	"github.com/mitchellh/mapstructure"
	"log"
	"time"
)

var RoleController = &roleController{}

type roleController struct {
}

func (r *roleController) InitRouter(router *net.Router) {
	g := router.Group("role")
	g.Use(middleware.Log())
	g.AddRouter("enterServer", r.enterServer)
	g.AddRouter("myProperty", r.myProperty, middleware.CheckRole())
	g.AddRouter("posTagList", r.posTagList)
	g.AddRouter("create", r.create)
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
	err = logic.RoleService.EnterServer(uid, rspObj, req)
	if err != nil {
		rspObj.Time = time.Now().UnixMilli()
		rsp.Body.Msg = rspObj
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
	respObj.RoleRes = (logic.RoleResService.GetRoleRes(rid).ToModel()).(model.RoleRes)
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

func (r *roleController) posTagList(req *net.WsMsgReq, resp *net.WsMsgResp) {
	respObj := &model.PosTagListRsp{}

	resp.Body.Seq = req.Body.Seq
	resp.Body.Name = req.Body.Name

	role, err := req.Conn.GetProperty("role")
	if err != nil {
		resp.Body.Code = constant.SessionInvalid
		return
	}
	rid := role.(data.Role).RId

	pts, err := logic.RoleAttrService.GetTagList(rid)
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	respObj.PosTags = pts
	resp.Body.Code = constant.OK
	resp.Body.Msg = respObj
}

func (r *roleController) create(req *net.WsMsgReq, resp *net.WsMsgResp) {
	reqObj := &model.CreateRoleReq{}
	respObj := &model.CreateRoleRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)

	resp.Body.Seq = req.Body.Seq
	resp.Body.Name = req.Body.Name
	role := &data.Role{}
	ok, err := db.Engine.Where("uid=?", reqObj.UId).Get(role)
	if err != nil {
		resp.Body.Code = constant.DBError
		return
	}
	if ok {
		resp.Body.Code = constant.RoleAlreadyCreate
		return
	}
	role.UId = reqObj.UId
	role.Sex = reqObj.Sex
	role.NickName = reqObj.NickName
	role.Balance = 0
	role.HeadId = reqObj.HeadId
	role.CreatedAt = time.Now()
	role.LoginTime = time.Now()
	_, err = db.Engine.InsertOne(role)
	if err != nil {
		resp.Body.Code = constant.DBError
		return
	}
	respObj.Role = role.ToModel().(model.Role)
	resp.Body.Code = constant.OK
	resp.Body.Msg = respObj
}
