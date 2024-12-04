package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"Turing-Go/utils"
	"log"
	"time"
)

var RoleService = &roleService{}

type roleService struct {
}

func (r *roleService) EnterServer(uid int, rsp *model.EnterServerRsp, conn net.WSConn) error {
	role := data.Role{}
	ok, err := db.Engine.Table(role).Where("uid=?", uid).Get(&role)
	if err != nil {
		log.Println("查询角色异常", err)
		return common.New(constant.DBError, "数据库错误")
	}
	if ok {
		rid := role.RId
		roleRes := &data.RoleRes{}
		ok, err = db.Engine.Table(roleRes).Where("rid=?", rid).Get(roleRes)
		if err != nil {
			log.Println("查询角色资源异常", err)
			return common.New(constant.DBError, "数据库错误")
		}
		if !ok {
			roleRes.RId = rid
			roleRes.Gold = gameConfig.Basic.Role.Gold
			roleRes.Decree = gameConfig.Basic.Role.Decree
			roleRes.Wood = gameConfig.Basic.Role.Wood
			roleRes.Iron = gameConfig.Basic.Role.Iron
			roleRes.Stone = gameConfig.Basic.Role.Stone
			roleRes.Grain = gameConfig.Basic.Role.Grain
			_, err = db.Engine.Table(roleRes).Insert(roleRes)
			if err != nil {
				log.Println("插入角色资源异常", err)
				return common.New(constant.DBError, "数据库错误")
			}
		}
		rsp.RoleRes = roleRes.ToModel().(model.RoleRes)
		rsp.Role = role.ToModel().(model.Role)
		rsp.Time = time.Now().UnixMilli()
		token, err := utils.Award(uid)
		if err != nil {
			log.Println("生成新鉴权异常", err)
			return common.New(constant.DBError, "数据库错误")
		}
		rsp.Token = token
		conn.SetProperty("role", role)
		if err := RoleAttrService.TryCreate(rid, conn); err != nil {
			return common.New(constant.DBError, "数据库错误")
		}
		// 初始化城池
		if err := RoleCityService.InitCity(rid, role.NickName, conn); err != nil {
			return common.New(constant.DBError, "数据库错误")
		}
	} else {
		return common.New(constant.RoleNotExist, "角色不存在")
	}
	return nil
}
