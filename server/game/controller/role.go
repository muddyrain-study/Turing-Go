package controller

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/net"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"Turing-Go/utils"
	"github.com/mitchellh/mapstructure"
	"log"
	"time"
)

var DefaultRoleHandler = &RoleHandler{}

type RoleHandler struct {
}

func (r *RoleHandler) InitRouter(router *net.Router) {
	g := router.Group("role")
	g.AddRouter("enterServer", r.enterServer)
}

func (r *RoleHandler) enterServer(req *net.WsMsgReq, rsp *net.WsMsgResp) {
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
	role := data.Role{}
	ok, err := db.Engine.Table(role).Where("uid=?", uid).Get(&role)
	if err != nil {
		log.Println("enterServer db.Engine.Table(role).Where error:", err)
		rsp.Body.Code = constant.DBError
		return
	}
	if !ok {
		rsp.Body.Code = constant.RoleNotExist
		return
	}
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj

	rid := role.RId
	roleRes := &data.RoleRes{}
	ok, err = db.Engine.Table(roleRes).Where("rid=?", rid).Get(roleRes)
	if err != nil {
		log.Println("enterServer db.Engine.Table(roleRes).Where error:", err)
		rsp.Body.Code = constant.DBError
		return
	}
	if !ok {
		roleRes.RId = rid
		roleRes.Gold = gameConfig.Basic.Role.Gold
		roleRes.Decree = gameConfig.Basic.Role.Decree
		roleRes.Wood = gameConfig.Basic.Role.Wood
		roleRes.Iron = gameConfig.Basic.Role.Iron
		roleRes.Stone = gameConfig.Basic.Role.Stone
		roleRes.Grain = gameConfig.Basic.Role.Grain
		_, err = db.Engine.Table(roleRes).Insert(roleRes)
		if err != nil {
			log.Println("enterServer db.Engine.Table(roleRes).Insert error:", err)
			rsp.Body.Code = constant.DBError
			return
		}
		return
	}
	rspObj.RoleRes = roleRes.ToModel().(model.RoleRes)
	rspObj.Role = role.ToModel().(model.Role)
	rspObj.Time = time.Now().UnixMilli()
	token, err := utils.Award(uid)
	if err != nil {
		log.Println("enterServer utils.Award error:", err)
		rsp.Body.Code = constant.DBError
		return
	}
	rspObj.Token = token

}
