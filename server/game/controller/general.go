package controller

import (
	"Turing-Go/constant"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/logic"
	"Turing-Go/server/game/middleware"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"github.com/mitchellh/mapstructure"
)

var GeneralController = &generalController{}

type generalController struct {
}

func (g *generalController) InitRouter(router *net.Router) {
	r := router.Group("general")
	r.Use(middleware.Log())
	r.AddRouter("myGenerals", g.myGenerals, middleware.CheckRole())
	r.AddRouter("drawGeneral", g.drawGeneral, middleware.CheckRole())
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

func (g *generalController) drawGeneral(req *net.WsMsgReq, resp *net.WsMsgResp) {
	//1. 计算抽卡花费的金钱
	//2. 判断金钱是否足够
	//3. 抽卡的次数 + 已有的武将 卡池是否足够
	//4. 随机生成武将即可（之前有实现）
	//5. 金币的扣除
	reqObj := &model.DrawGeneralReq{}
	err := mapstructure.Decode(req.Body.Msg, reqObj)
	if err != nil {
		resp.Body.Code = constant.InvalidParam
		return
	}
	resp.Body.Seq = req.Body.Seq
	resp.Body.Name = req.Body.Name
	respObj := &model.DrawGeneralRsp{}
	resp.Body.Msg = respObj
	resp.Body.Code = constant.OK
	role, _ := req.Conn.GetProperty("role")
	rid := role.(data.Role).RId
	cost := gameConfig.Basic.General.DrawGeneralCost * reqObj.DrawTimes
	if !logic.RoleResService.IsEnoughGold(rid, cost) {
		resp.Body.Code = constant.GoldNotEnough
		return
	}
	limit := gameConfig.Basic.General.Limit

	gs, err := logic.GeneralService.GetGenerals(rid)
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	if (len(gs) + reqObj.DrawTimes) > limit {
		resp.Body.Code = constant.OutGeneralLimit
		return
	}
	mgs := logic.GeneralService.Draw(rid, reqObj.DrawTimes)
	logic.RoleResService.CostGold(rid, cost)
	respObj.Generals = mgs
}
