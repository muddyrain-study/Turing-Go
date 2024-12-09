package controller

import (
	"Turing-Go/constant"
	"Turing-Go/net"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/logic"
	"Turing-Go/server/game/middleware"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"time"
)

var InteriorController = &interiorController{}

type interiorController struct {
}

func (r *interiorController) InitRouter(router *net.Router) {
	g := router.Group("interior")
	g.Use(middleware.Log())
	g.AddRouter("openCollect", r.openCollect, middleware.CheckRole())
	g.AddRouter("collect", r.collect, middleware.CheckRole())
}

func (r *interiorController) openCollect(req *net.WsMsgReq, resp *net.WsMsgResp) {
	respObj := &model.OpenCollectionRsp{}
	resp.Body.Seq = req.Body.Seq
	resp.Body.Name = req.Body.Name
	role, _ := req.Conn.GetProperty("role")
	rid := role.(data.Role).RId
	ra := logic.RoleAttrService.Get(rid)
	if ra != nil {
		// 次数
		respObj.CurTimes = ra.CollectTimes
		respObj.Limit = gameConfig.Basic.Role.CollectTimesLimit
		// 间隔时间
		interval := gameConfig.Basic.Role.CollectInterval
		if ra.LastCollectTime.IsZero() {
			respObj.NextTime = 0
		} else {
			if respObj.CurTimes >= respObj.Limit {
				y, m, d := ra.LastCollectTime.Add(24 * time.Hour).Date()
				ti := time.Date(y, m, d, 0, 0, 0, 0, time.FixedZone("CST", 8*3600))
				respObj.NextTime = ti.UnixMilli()
			} else {
				ti := ra.LastCollectTime.Add(time.Duration(interval) * time.Second)
				respObj.NextTime = ti.UnixMilli()
			}
		}
	}
	resp.Body.Code = constant.OK
	resp.Body.Msg = respObj
}

func (r *interiorController) collect(req *net.WsMsgReq, resp *net.WsMsgResp) {
	respObj := &model.CollectionRsp{}

	resp.Body.Msg = respObj
	resp.Body.Code = constant.OK

	role, _ := req.Conn.GetProperty("role")
	rid := role.(data.Role).RId
	ra := logic.RoleAttrService.Get(rid)
	if ra == nil {
		resp.Body.Code = constant.DBError
		return
	}
	rs := logic.RoleResService.GetRoleRes(rid)
	if rs == nil {
		resp.Body.Code = constant.DBError
		return
	}
	yield := logic.RoleResService.GetYield(rid)
	rs.Gold += yield.Gold
	rs.SyncExecute()
	respObj.Gold = yield.Gold
	curTime := time.Now()
	limit := gameConfig.Basic.Role.CollectTimesLimit
	interval := gameConfig.Basic.Role.CollectInterval
	lastTime := ra.LastCollectTime
	if curTime.YearDay() != lastTime.YearDay() || curTime.Year() != lastTime.Year() {
		ra.CollectTimes = 0
		ra.LastCollectTime = time.Time{}
	}
	ra.CollectTimes += 1
	ra.LastCollectTime = curTime
	ra.SyncExecute()
	respObj.Limit = limit
	respObj.CurTimes = ra.CollectTimes
	if respObj.CurTimes >= respObj.Limit {
		y, m, d := ra.LastCollectTime.Add(24 * time.Hour).Date()
		ti := time.Date(y, m, d, 0, 0, 0, 0, time.FixedZone("CST", 8*3600))
		respObj.NextTime = ti.UnixMilli()
	} else {
		ti := ra.LastCollectTime.Add(time.Duration(interval) * time.Second)
		respObj.NextTime = ti.UnixMilli()
	}
}
