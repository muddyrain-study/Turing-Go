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

func armyIsInView(rid, x, y int) bool {
	//简单点 先设为true
	return true
}
