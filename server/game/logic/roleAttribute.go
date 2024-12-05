package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"encoding/json"
	"log"
	"sync"
)

var RoleAttrService = &roleAttrService{
	attrs: make(map[int]*data.RoleAttribute),
}

type roleAttrService struct {
	attrs map[int]*data.RoleAttribute
	mutex sync.RWMutex
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
		roleAttribute.PosTags = ""
		_, err := db.Engine.Table(roleAttribute).Insert(&roleAttribute)
		if err != nil {
			log.Println("插入角色属性异常", err)
			return common.New(constant.DBError, "数据库错误")
		}
	}
	r.mutex.Lock()
	r.attrs[rid] = &roleAttribute
	r.mutex.Unlock()
	return nil

}

func (r *roleAttrService) GetTagList(rid int) ([]model.PosTag, error) {
	var err error
	roleAttribute, ok := r.attrs[rid]
	if !ok {
		roleAttribute = &data.RoleAttribute{}
		ok, err = db.Engine.Table(roleAttribute).Where("rid=?", rid).Get(roleAttribute)
		if err != nil {
			log.Println("查询角色属性异常", err)
			return nil, common.New(constant.DBError, "数据库错误")
		}
	}
	posTags := make([]model.PosTag, 0)
	if !ok {
		return posTags, nil
	}
	tags := roleAttribute.PosTags
	if tags != "" {
		err := json.Unmarshal([]byte(tags), &posTags)
		if err != nil {
			log.Println("解析标签异常", err)
			return nil, common.New(constant.DBError, "数据库错误")
		}
	}
	return posTags, nil
}
