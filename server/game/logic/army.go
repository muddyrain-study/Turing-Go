package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/server/common"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"log"
)

var ArmyService = &armyService{}

type armyService struct {
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

func (g *armyService) ScanBlock(obj *model.ScanBlockReq) ([]model.Army, error) {
	return nil, nil
}
