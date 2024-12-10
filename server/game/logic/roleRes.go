package logic

import (
	"Turing-Go/db"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/model/data"
	"log"
)

var RoleResService = &roleResService{}

type roleResService struct {
}

func (r *roleResService) GetRoleRes(rid int) *data.RoleRes {
	roleRes := &data.RoleRes{}
	ok, err := db.Engine.Table(roleRes).Where("rid=?", rid).Get(roleRes)
	if err != nil {
		log.Println("查询角色资源异常", err)
		return roleRes
	}
	if !ok {
		return roleRes
	}
	return roleRes
}

func (r *roleResService) GetYield(rid int) data.Yield {
	//产量 建筑 城市固定收益 = 最终的产量
	rbYield := RoleBuildService.GetYield(rid)
	rcYield := CityFacilityService.GetYield(rid)
	var y data.Yield

	y.Gold = rbYield.Gold + rcYield.Gold + gameConfig.Basic.Role.GoldYield
	y.Stone = rbYield.Stone + rcYield.Stone + gameConfig.Basic.Role.StoneYield
	y.Iron = rbYield.Iron + rcYield.Iron + gameConfig.Basic.Role.IronYield
	y.Grain = rbYield.Grain + rcYield.Grain + gameConfig.Basic.Role.GrainYield
	y.Wood = rbYield.Wood + rcYield.Wood + gameConfig.Basic.Role.WoodYield

	return y
}

func (r *roleResService) IsEnoughGold(rid int, cost int) bool {
	rr := r.GetRoleRes(rid)
	return rr.Gold >= cost
}

func (r *roleResService) CostGold(rid int, cost int) {
	rr := r.GetRoleRes(rid)
	if rr.Gold >= cost {
		rr.Gold -= cost
		rr.SyncExecute()
	}
}
