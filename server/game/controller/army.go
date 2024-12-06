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
	"log"
)

var ArmyController = &armyController{}

type armyController struct {
}

func (g *armyController) InitRouter(router *net.Router) {
	r := router.Group("army")
	r.Use(middleware.Log())
	r.AddRouter("myList", g.myList, middleware.CheckRole())
}

func (g *armyController) myList(req *net.WsMsgReq, resp *net.WsMsgResp) {
	reqObj := &model.ArmyListReq{}
	err := mapstructure.Decode(req.Body.Msg, reqObj)
	if err != nil {
		log.Println("解析军队列表请求失败: ", err)
		return
	}
	respObj := &model.ArmyListRsp{}
	resp.Body.Msg = respObj
	resp.Body.Code = constant.OK
	resp.Body.Seq = req.Body.Seq
	resp.Body.Name = req.Body.Name
	role, _ := req.Conn.GetProperty("role")
	r := role.(data.Role)
	arms, err := logic.ArmyService.GetArmiesByCity(r.RId, reqObj.CityId)
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	respObj.Armys = arms
	respObj.CityId = reqObj.CityId
	resp.Body.Msg = respObj
}
