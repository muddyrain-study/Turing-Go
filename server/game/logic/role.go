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

func (r *roleService) EnterServer(uid int, rsp *model.EnterServerRsp, req *net.WsMsgReq) error {
	role := data.Role{}
	session := db.Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		log.Println("事务出错", err)
		return common.New(constant.DBError, "查询角色出错")
	}
	req.Context.Set("session", session)
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
			_, err = session.Table(roleRes).Insert(roleRes)
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
		req.Conn.SetProperty("role", role)
		if err := RoleAttrService.TryCreate(rid, req); err != nil {
			err := session.Rollback()
			if err != nil {
				log.Println("事务回滚出错", err)
			}
			return common.New(constant.DBError, "数据库错误")
		}
		// 初始化城池
		if err := RoleCityService.InitCity(rid, role.NickName, req); err != nil {
			err := session.Rollback()
			if err != nil {
				log.Println("事务回滚出错", err)
			}
			return common.New(constant.DBError, "数据库错误")
		}
	} else {
		return common.New(constant.RoleNotExist, "角色不存在")
	}
	err = session.Commit()
	if err != nil {
		log.Println("事务提交提交出错", err)
		return common.New(constant.DBError, "数据库错误")
	}
	return nil
}

func (r *roleService) GetRoleRes(rid int) (model.RoleRes, error) {
	roleRes := &data.RoleRes{}
	ok, err := db.Engine.Table(roleRes).Where("rid=?", rid).Get(roleRes)
	if err != nil {
		log.Println("查询角色资源异常", err)
		return model.RoleRes{}, common.New(constant.DBError, "数据库错误")
	}
	if !ok {
		return model.RoleRes{}, common.New(constant.DBError, "资源不存在")
	}
	return roleRes.ToModel().(model.RoleRes), nil
}
