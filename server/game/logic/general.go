package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/server/common"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"log"
)

var GeneralService = &generalService{}

type generalService struct {
}

func (g *generalService) GetGenerals(rid int) ([]model.General, error) {
	mrs := make([]data.General, 0)
	mr := &data.General{}
	err := db.Engine.Table(mr).Where("rid=?", rid).Find(&mrs)
	if err != nil {
		log.Println("武将查询出错", err)
		return nil, common.New(constant.DBError, "武将查询出错")
	}
	modelMrs := make([]model.General, len(mrs))
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.General))
	}
	return modelMrs, nil
}
