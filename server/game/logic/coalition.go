package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/server/common"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"log"
	"sync"
)

var CoalitionService = &coalitionService{
	unions: make(map[int]*data.Coalition),
}

type coalitionService struct {
	mutex  sync.RWMutex
	unions map[int]*data.Coalition
}

func (c *coalitionService) Load() {
	rr := make([]*data.Coalition, 0)
	err := db.Engine.Table(new(data.Coalition)).Where("state=?", data.UnionRunning).Find(&rr)
	if err != nil {
		log.Println("coalitionService load error", err)
	}
	for _, v := range rr {
		c.unions[v.Id] = v
	}
}

func (c *coalitionService) List() ([]model.Union, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	uns := make([]model.Union, 0)
	for _, v := range c.unions {
		mas := make([]model.Major, 0)
		role, err := RoleService.Get(v.Chairman)
		if err != nil {
			log.Println("盟主角色不存在", err)
		}
		if role != nil {
			ma := model.Major{
				RId:   role.RId,
				Name:  role.NickName,
				Title: model.UnionChairman,
			}
			mas = append(mas, ma)
		}
		role2, err := RoleService.Get(v.ViceChairman)
		if err != nil {
			log.Println("副盟主角色不存在", err)
		}
		if role2 != nil {
			ma := model.Major{
				RId:   role2.RId,
				Name:  role2.NickName,
				Title: model.UnionViceChairman,
			}
			mas = append(mas, ma)
		}

		union := v.ToModal().(model.Union)
		union.Major = mas
		uns = append(uns, union)
	}
	return uns, nil
}

func (c *coalitionService) ListCoalition() []*data.Coalition {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	coalitions := make([]*data.Coalition, 0)
	for _, v := range c.unions {
		coalitions = append(coalitions, v)
	}
	return coalitions
}

func (c *coalitionService) Get(id int) (model.Union, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	coa, ok := c.unions[id]
	if ok {
		union := coa.ToModal().(model.Union)
		mas := make([]model.Major, 0)
		role, err := RoleService.Get(coa.Chairman)
		if err != nil {
			log.Println("盟主角色不存在", err)
		}
		if role != nil {
			ma := model.Major{
				RId:   role.RId,
				Name:  role.NickName,
				Title: model.UnionChairman,
			}
			mas = append(mas, ma)
		}
		role2, err := RoleService.Get(coa.ViceChairman)
		if err != nil {
			log.Println("副盟主角色不存在", err)
		}
		if role2 != nil {
			ma := model.Major{
				RId:   role2.RId,
				Name:  role2.NickName,
				Title: model.UnionViceChairman,
			}
			mas = append(mas, ma)
		}
		union.Major = mas
		return union, nil
	}
	return model.Union{}, nil
}

func (c *coalitionService) GetCoalition(id int) *data.Coalition {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	coa, ok := c.unions[id]
	if ok {
		return coa
	}
	return nil
}

func (c *coalitionService) GetListApply(unionId int, state int) ([]model.ApplyItem, error) {
	applys := make([]data.CoalitionApply, 0)
	err := db.Engine.Table(new(data.CoalitionApply)).
		Where("union_id=? and state=?", unionId, state).
		Find(&applys)
	if err != nil {
		log.Println("coalitionService GetListApply find error", err)
		return nil, common.New(constant.DBError, "数据库错误")
	}
	ais := make([]model.ApplyItem, 0)
	for _, v := range applys {
		var ai model.ApplyItem
		ai.Id = v.Id
		role, err := RoleService.Get(v.RId)
		if err != nil {
			log.Println("查询角色异常", err)
			return nil, common.New(constant.DBError, "数据库错误")
		}
		ai.NickName = role.NickName
		ai.RId = role.RId
		ais = append(ais, ai)
	}
	return ais, nil
}
