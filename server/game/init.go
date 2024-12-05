package game

import (
	"Turing-Go/db"
	"Turing-Go/net"
	"Turing-Go/server/game/controller"
	"Turing-Go/server/game/gameConfig"
)

var Router = &net.Router{}

func Init() {
	db.TestDB()
	gameConfig.Basic.Load()
	gameConfig.MapBuildConf.Load()
	gameConfig.MapRes.Load()
	gameConfig.FacilityConf.Load()
	initRouter()
}

func initRouter() {
	controller.RoleController.InitRouter(Router)
	controller.NationMapController.InitRouter(Router)
}
