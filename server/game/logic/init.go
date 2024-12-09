package logic

import "Turing-Go/server/game/model/data"

func BeforeInit() {
	data.GetYield = RoleResService.GetYield
}
