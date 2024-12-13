package logic

import (
	"Turing-Go/db"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/model/data"
	"log"
	"time"
)

var RoleResService = &roleResService{
	rolesRes: make(map[int]*data.RoleRes),
}

type roleResService struct {
	rolesRes map[int]*data.RoleRes
}

func (r *roleResService) Load() {
	rr := make([]*data.RoleRes, 0)
	err := db.Engine.Find(&rr)
	if err != nil {
		log.Println("加载角色资源异常", err)
	}
	for _, v := range rr {
		r.rolesRes[v.RId] = v
	}
	go r.produce()
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

func (r *roleResService) TryUseNeed(rid int, need gameConfig.NeedRes) bool {
	rr := r.GetRoleRes(rid)
	if need.Decree <= rr.Decree && need.Grain <= rr.Grain &&
		need.Stone <= rr.Stone && need.Wood <= rr.Wood &&
		need.Iron <= rr.Iron && need.Gold <= rr.Gold {
		rr.Decree -= need.Decree
		rr.Iron -= need.Iron
		rr.Wood -= need.Wood
		rr.Stone -= need.Stone
		rr.Grain -= need.Grain
		rr.Gold -= need.Gold
		rr.SyncExecute()
		return true
	} else {
		return false
	}
}

func (r *roleResService) produce() {
	for {
		// 一直获取产量 隔一段时间获取
		recoveryTime := gameConfig.Basic.Role.RecoveryTime
		time.Sleep(time.Duration(recoveryTime) * time.Second)
		var index int
		for _, v := range r.rolesRes {
			capacity := r.getDepotCapacity(v.RId)
			yield := r.GetYield(v.RId)
			if v.Wood < capacity {
				v.Wood += yield.Wood / 6
			}
			if v.Iron < capacity {
				v.Iron += yield.Iron / 6
			}
			if v.Stone < capacity {
				v.Stone += yield.Stone / 6
			}
			if v.Grain < capacity {
				v.Grain += yield.Grain / 6
			}
			if index%6 == 0 {
				if v.Decree < gameConfig.Basic.Role.Decree {
					v.Decree += 1
				}
			}
			v.SyncExecute()
		}
	}
}

func (r *roleResService) getDepotCapacity(rid int) int {
	return CityFacilityService.GetCapacity(rid) + gameConfig.Basic.Role.DepotCapacity
}
