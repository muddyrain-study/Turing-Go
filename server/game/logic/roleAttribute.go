package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/model/data"
	"log"
)

var RoleAttrService = &roleAttrService{}

type roleAttrService struct {
}

func (r *roleAttrService) TryCreate(rid int, conn net.WSConn) error {
	roleAttribute := data.RoleAttribute{}
	ok, err := db.Engine.Table(roleAttribute).Where("rid=?", rid).Get(&roleAttribute)
	if err != nil {
		log.Println("查询角色属性异常", err)
		return common.New(constant.DBError, "数据库错误")
	}
	if !ok {
		roleAttribute.RId = rid
		roleAttribute.UnionId = 0
		roleAttribute.ParentId = 0
		_, err := db.Engine.Table(roleAttribute).Insert(&roleAttribute)
		if err != nil {
			log.Println("插入角色属性异常", err)
			return common.New(constant.DBError, "数据库错误")
		}
	}
	return nil

}
