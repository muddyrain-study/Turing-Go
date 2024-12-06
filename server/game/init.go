package game

import (
	"Turing-Go/db"
	"Turing-Go/net"
	"Turing-Go/server/game/controller"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/gameConfig/general"
	"Turing-Go/server/game/logic"
)

var Router = &net.Router{}

func Init() {
	db.TestDB()
	gameConfig.Basic.Load()
	gameConfig.MapBuildConf.Load()
	gameConfig.MapRes.Load()
	gameConfig.FacilityConf.Load()
	gameConfig.Skill.Load()
	general.General.Load()
	logic.RoleBuildService.Load()
	initRouter()
}

func initRouter() {
	controller.RoleController.InitRouter(Router)
	controller.NationMapController.InitRouter(Router)
	controller.GeneralController.InitRouter(Router)
	controller.ArmyController.InitRouter(Router)
	controller.WarController.InitRouter(Router)
	controller.SkillController.InitRouter(Router)
}
