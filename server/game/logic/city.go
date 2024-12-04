package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/global"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"log"
	"math/rand"
	"time"
)

var RoleCityService = &roleCityService{}

type roleCityService struct {
}

func (r *roleCityService) InitCity(rid int, nickName string, conn net.WSConn) error {
	roleCity := data.MapRoleCity{}
	ok, err := db.Engine.Table(roleCity).Where("rid=?", rid).Get(&roleCity)
	if err != nil {
		log.Println("查询角色城池异常", err)
		return common.New(constant.DBError, "数据库错误")
	}
	if !ok {
		roleCity.X = rand.Intn(global.MapWidth)
		roleCity.Y = rand.Intn(global.MapHeight)

		roleCity.RId = rid
		roleCity.Name = nickName
		roleCity.CurDurable = gameConfig.Basic.City.Durable
		roleCity.CreatedAt = time.Now()
		roleCity.IsMain = 1
		_, err := db.Engine.Table(roleCity).Insert(&roleCity)
		if err != nil {
			log.Println("插入角色城池异常", err)
			return common.New(constant.DBError, "数据库错误")
		}
	}
	return nil

}

func (r *roleCityService) GetRoleCities(rid int) ([]model.MapRoleCity, error) {
	cities := make([]data.MapRoleCity, 0)
	city := &data.MapRoleCity{}
	err := db.Engine.Table(city).Where("rid=?", rid).Find(&cities)
	modelCities := make([]model.MapRoleCity, len(cities))

	if err != nil {
		log.Println("查询角色城池异常", err)
		return modelCities, common.New(constant.DBError, "数据库错误")
	}
	for i, roleCity := range cities {
		modelCities[i] = roleCity.ToModel().(model.MapRoleCity)
	}
	return modelCities, nil
}
