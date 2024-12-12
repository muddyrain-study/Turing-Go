package controller

import (
	"Turing-Go/constant"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/gameConfig/general"
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
	r.AddRouter("dispose", g.dispose, middleware.CheckRole())
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

func (g *armyController) dispose(req *net.WsMsgReq, resp *net.WsMsgResp) {
	reqObj := &model.DisposeReq{}
	respObj := &model.DisposeRsp{}
	err := mapstructure.Decode(req.Body.Msg, reqObj)
	if err != nil {
		log.Println("解析军队处理请求失败: ", err)
		resp.Body.Code = constant.InvalidParam
		return
	}
	resp.Body.Code = constant.OK
	resp.Body.Seq = req.Body.Seq
	resp.Body.Name = req.Body.Name
	resp.Body.Msg = respObj

	r, _ := req.Conn.GetProperty("role")
	role := r.(data.Role)

	if reqObj.Order <= 0 || reqObj.Order > 5 || reqObj.Position < -1 || reqObj.Position > 2 {
		resp.Body.Code = constant.InvalidParam
		return
	}

	city, ok := logic.RoleCityService.Get(reqObj.CityId)
	if !ok {
		resp.Body.Code = constant.CityNotExist
		return
	}
	if role.RId != city.RId {
		resp.Body.Code = constant.CityNotMe
		return
	}
	//校场每升一级一个队伍
	level := logic.CityFacilityService.GetFacilityLevel(reqObj.CityId, gameConfig.JiaoChang)
	if level <= 0 || reqObj.Order > level {
		resp.Body.Code = constant.ArmyNotEnough
		return
	}
	newGen, ok := logic.GeneralService.Get(reqObj.GeneralId)
	if !ok {
		resp.Body.Code = constant.GeneralNotFound
		return
	}
	if newGen.RId != role.RId {
		resp.Body.Code = constant.GeneralNotMe
		return
	}
	army, ok := logic.ArmyService.GetOrCreate(role.RId, reqObj.CityId, reqObj.Order)

	if (army.FromX > 0 && army.FromX != city.X) || (army.FromY > 0 && army.FromY != city.Y) {
		resp.Body.Code = constant.ArmyIsOutside
		return
	}
	if reqObj.Position == -1 {
		// 下阵
		for pos, g := range army.Gens {
			if g != nil && g.Id == reqObj.GeneralId {
				//征兵中不能下阵
				if army.PositionCanModify(pos) == false {
					if army.Cmd == data.ArmyCmdConscript {
						resp.Body.Code = constant.GeneralBusy
					} else {
						resp.Body.Code = constant.ArmyBusy
					}
					return
				}
				army.GeneralArray[pos] = 0
				army.SoldierArray[pos] = 0
				army.Gens[pos] = nil
				army.SyncExecute()
				break
			}
		}
		newGen.Order = 0
		newGen.CityId = 0
		newGen.SyncExecute()
	} else {
		// 上阵
		//征兵中不能上阵
		if army.PositionCanModify(reqObj.Position) == false {
			if army.Cmd == data.ArmyCmdConscript {
				resp.Body.Code = constant.GeneralBusy
			} else {
				resp.Body.Code = constant.ArmyBusy
			}
			return
		}
		if newGen.CityId != 0 {
			resp.Body.Code = constant.GeneralBusy
			return
		}

		if logic.ArmyService.IsRepeat(role.RId, newGen.CfgId) == false {
			resp.Body.Code = constant.GeneralRepeat
			return
		}

		//判断是否能配前锋
		level := logic.CityFacilityService.GetFacilityLevel(city.CityId, gameConfig.TongShuaiTing)
		if reqObj.Position == 2 && (ok == false || level < reqObj.Order) {
			resp.Body.Code = constant.TongShuaiNotEnough
			return
		}

		//判断cost
		cost := general.General.Cost(newGen.CfgId)
		for i, g := range army.Gens {
			if g == nil || i == reqObj.Position {
				continue
			}
			cost += general.General.Cost(g.CfgId)
		}

		if logic.RoleCityService.GetCityCost(city.CityId) < cost {
			resp.Body.Code = constant.CostNotEnough
			return
		}

		oldG := army.Gens[reqObj.Position]
		if oldG != nil {
			//旧的下阵
			oldG.CityId = 0
			oldG.Order = 0
			oldG.SyncExecute()
		}

		army.GeneralArray[reqObj.Position] = reqObj.GeneralId
		army.SoldierArray[reqObj.Position] = 0
		army.Gens[reqObj.Position] = newGen

		newGen.Order = reqObj.Order
		newGen.CityId = reqObj.CityId
		newGen.SyncExecute()
	}
	army.FromX = city.X
	army.FromY = city.Y
	army.SyncExecute()
	//队伍
	respObj.Army = army.ToModel().(model.Army)
}
