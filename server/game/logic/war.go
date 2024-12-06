package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/server/common"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"log"
)

type warService struct{}

var WarService = &warService{}

func (w *warService) GetWarReports(rid int) ([]model.WarReport, error) {
	mrs := make([]data.WarReport, 0)
	mr := &data.WarReport{}
	err := db.Engine.Table(mr).Where("a_rid=? or d_rid=?", rid, rid).Limit(30, 0).Desc("cTime").Find(&mrs)
	if err != nil {
		log.Println("战报查询出错", err)
		return nil, common.New(constant.DBError, "战报查询出错")
	}
	modelMrs := make([]model.WarReport, 0)
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.WarReport))
	}
	return modelMrs, nil
}
