package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/server/common"
	"Turing-Go/server/game/global"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"Turing-Go/utils"
	"log"
	"sync"
)

var ArmyService = &armyService{
	passBy:         sync.RWMutex{},
	passByPosArmys: make(map[int]map[int]*data.Army),
}

type armyService struct {
	passBy         sync.RWMutex
	passByPosArmys map[int]map[int]*data.Army //玩家路过位置的军队 key:posId,armyId
}

func (g *armyService) GetArmies(rid int) ([]model.Army, error) {
	mrs := make([]data.Army, 0)
	mr := &data.Army{}
	err := db.Engine.Table(mr).Where("rid=?", rid).Find(&mrs)
	if err != nil {
		log.Println("军队查询出错", err)
		return nil, common.New(constant.DBError, "军队查询出错")
	}
	modelMrs := make([]model.Army, 0)
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.Army))
	}
	return modelMrs, nil
}

func (g *armyService) GetDbArmies(rid int) ([]*data.Army, error) {
	mrs := make([]*data.Army, 0)
	mr := &data.Army{}
	err := db.Engine.Table(mr).Where("rid=?", rid).Find(&mrs)
	if err != nil {
		log.Println("军队查询出错", err)
		return nil, common.New(constant.DBError, "军队查询出错")
	}
	modelMrs := make([]*data.Army, 0)
	for _, v := range mrs {
		modelMrs = append(modelMrs, v)
	}
	return modelMrs, nil
}

func (g *armyService) GetArmiesByCity(rid int, cityId int) ([]model.Army, error) {
	mrs := make([]data.Army, 0)
	mr := &data.Army{}
	err := db.Engine.Table(mr).Where("rid=?&cityId=?", rid, cityId).Find(&mrs)
	if err != nil {
		log.Println("军队查询出错", err)
		return nil, common.New(constant.DBError, "军队查询出错")
	}
	modelMrs := make([]model.Army, 0)
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.Army))
	}
	return modelMrs, nil
}

func (g *armyService) ScanBlock(roleId int, req *model.ScanBlockReq) ([]model.Army, error) {
	x := req.X
	y := req.Y
	length := req.Length
	if x < 0 || x >= global.MapWidth || y < 0 || y >= global.MapHeight {
		return nil, common.New(constant.InvalidParam, "坐标错误")
	}

	maxX := utils.MinInt(global.MapWidth, x+length-1)
	maxY := utils.MinInt(global.MapHeight, y+length-1)
	out := make([]model.Army, 0)

	g.passBy.RLock()
	for i := x; i <= maxX; i++ {
		for j := y; j <= maxY; j++ {

			posId := global.ToPosition(i, j)
			armys, ok := g.passByPosArmys[posId]
			if ok {
				//是否在视野范围内
				is := armyIsInView(roleId, i, j)
				if is == false {
					continue
				}
				for _, army := range armys {
					out = append(out, army.ToModel().(model.Army))
				}
			}
		}
	}
	g.passBy.RUnlock()
	return out, nil
}
func (g *armyService) updateGenerals(armys ...*data.Army) {
	for _, army := range armys {
		army.Gens = make([]*data.General, 0)
		for _, gid := range army.GeneralArray {
			if gid == 0 {
				army.Gens = append(army.Gens, nil)
			} else {
				g, _ := GeneralService.Get(gid)
				army.Gens = append(army.Gens, g)
			}
		}
	}
}
func (g *armyService) GetArmiesByCityAndOrder(rid int, cityId int, order int8) (*data.Army, bool) {
	mr := &data.Army{}
	ok, err := db.Engine.Table(mr).Where("rid=? AND cityId=? AND `order`=?", rid, cityId, order).Get(mr)
	if err != nil {
		log.Println("军队查询出错", err)
		return nil, false
	}
	if ok {
		g.updateGenerals(mr)
		return mr, true
	} else {
		return nil, false
	}
}

func (g *armyService) GetOrCreate(rid int, cid int, order int8) (*data.Army, bool) {
	army, ok := g.GetArmiesByCityAndOrder(rid, cid, order)
	if ok {
		return army, true
	}
	//需要创建
	army = &data.Army{
		RId:                rid,
		Order:              order,
		CityId:             cid,
		Generals:           `[0,0,0]`,
		Soldiers:           `[0,0,0]`,
		GeneralArray:       []int{0, 0, 0},
		SoldierArray:       []int{0, 0, 0},
		ConscriptCnts:      `[0,0,0]`,
		ConscriptTimes:     `[0,0,0]`,
		ConscriptCntArray:  []int{0, 0, 0},
		ConscriptTimeArray: []int64{0, 0, 0},
	}
	g.updateGenerals(army)
	_, err := db.Engine.Table(army).Insert(army)
	if err == nil {
		return army, true
	} else {
		log.Println("armyService GetCreate err", err)
		return nil, false
	}
}

func (g *armyService) IsRepeat(rid int, cfgId int) bool {
	armys, err := g.GetDbArmies(rid)
	if err != nil {
		return true
	}
	for _, army := range armys {
		for _, g := range army.Gens {
			if g != nil {
				if g.CfgId == cfgId && g.CityId != 0 {
					return false
				}
			}
		}
	}
	return true
}

func armyIsInView(rid, x, y int) bool {
	//简单点 先设为true
	return true
}
