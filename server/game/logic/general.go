package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/server/common"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/gameConfig/general"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"encoding/json"
	"log"
	"time"
)

var GeneralService = &generalService{}

type generalService struct {
}

func (g *generalService) GetGenerals(rid int) ([]model.General, error) {
	mrs := make([]*data.General, 0)
	mr := &data.General{}
	err := db.Engine.Table(mr).Where("rid=?", rid).Find(&mrs)
	if err != nil {
		log.Println("武将查询出错", err)
		return nil, common.New(constant.DBError, "武将查询出错")
	}
	if len(mrs) <= 0 {
		var count = 0
		for {
			if count >= 3 {
				break
			}
			cfgId := general.General.Rand()
			if cfgId != 0 {
				gen, err := g.newGeneral(cfgId, rid, 1)
				if err != nil {
					log.Println("武将生成出错", err)
					continue
				}
				mrs = append(mrs, gen)
				count++
			}
		}
	}
	modelMrs := make([]model.General, 0)
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.General))
	}
	return modelMrs, nil
}
func (g *generalService) Draw(rid, nums int) []model.General {
	mrs := make([]*data.General, 0)
	for i := 0; i < nums; i++ {
		cfgId := general.General.Rand()
		if cfgId != 0 {
			gen, err := g.newGeneral(cfgId, rid, 0)
			if err != nil {
				log.Println("武将生成出错", err)
				continue
			}
			mrs = append(mrs, gen)
		}
	}
	modelMrs := make([]model.General, 0)
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.General))
	}
	return modelMrs
}

const (
	GeneralNormal      = 0 //正常
	GeneralComposeStar = 1 //星级合成
	GeneralConvert     = 2 //转换
)

func (g *generalService) newGeneral(cfgId int, rid int, level int8) (*data.General, interface{}) {
	cfg := general.General.GMap[cfgId]
	sa := make([]*model.GSkill, 3)
	ss, _ := json.Marshal(sa)
	gen := &data.General{
		PhysicalPower: gameConfig.Basic.General.PhysicalPowerLimit,
		RId:           rid,
		CfgId:         cfg.CfgId,
		Order:         0,
		CityId:        0,
		Level:         level,
		CreatedAt:     time.Now(),
		CurArms:       cfg.Arms[0],
		HasPrPoint:    0,
		UsePrPoint:    0,
		AttackDis:     0,
		ForceAdded:    0,
		StrategyAdded: 0,
		DefenseAdded:  0,
		SpeedAdded:    0,
		DestroyAdded:  0,
		Star:          cfg.Star,
		StarLv:        0,
		ParentId:      0,
		SkillsArray:   sa,
		Skills:        string(ss),
		State:         GeneralNormal,
	}
	_, err := db.Engine.Table(gen).Insert(gen)
	if err != nil {
		log.Println("武将插入出错", err)
		return nil, common.New(constant.DBError, "武将插入出错")
	}
	return gen, nil
}

func (g *generalService) Get(id int) (*data.General, bool) {
	gen := &data.General{}
	ok, err := db.Engine.Table(new(data.General)).Where("id=? and state=?", id, data.GeneralNormal).Get(gen)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	if ok {
		return gen, true
	}
	return nil, false
}
