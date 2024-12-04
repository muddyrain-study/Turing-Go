package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/server/common"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"log"
)

var RoleBuildService = &roleBuildService{}

type roleBuildService struct {
}

func (r *roleBuildService) GetBuilds(rid int, ) ([]model.MapRoleBuild, error) {
	mrs := make([]data.MapRoleBuild, 0)
	mr := &data.MapRoleBuild{}
	err := db.Engine.Table(mr).Where("rid=?", rid).Find(&mrs)
	modelBuilds := make([]model.MapRoleBuild, len(mrs))
	if err != nil {
		log.Println("查询角色建筑异常", err)
		return modelBuilds, common.New(constant.DBError, "数据库错误")
	}
	for i, roleBuild := range mrs {
		modelBuilds[i] = roleBuild.ToModel().(model.MapRoleBuild)
	}
	return modelBuilds, nil
}
