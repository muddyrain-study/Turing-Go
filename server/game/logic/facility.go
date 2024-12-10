package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/server/common"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/model/data"
	"encoding/json"
	"github.com/go-xorm/xorm"
	"log"
)

var CityFacilityService = &cityFacilityService{}

type cityFacilityService struct {
}

func (c *cityFacilityService) TryCreate(cid, rid int, session *xorm.Session) error {
	cf := &data.CityFacility{}
	ok, err := db.Engine.Table(cf).Where("cityId=?", cid).Get(cf)
	if err != nil {
		log.Println(err)
		return common.New(constant.DBError, "数据库错误")
	}
	if ok {
		return nil
	}
	list := gameConfig.FacilityConf.List
	fs := make([]data.Facility, len(list))
	for index, v := range list {
		fac := data.Facility{
			Name:         v.Name,
			PrivateLevel: 0,
			Type:         v.Type,
			UpTime:       0,
		}
		fs[index] = fac
	}
	fsJson, err := json.Marshal(fs)
	if err != nil {
		return common.New(constant.DBError, "转json出错")
	}
	cf.RId = rid
	cf.CityId = cid
	cf.Facilities = string(fsJson)
	if session != nil {
		_, err = session.Table(cf).Insert(cf)
	} else {
		_, err = db.Engine.Table(cf).Insert(cf)
	}
	if err != nil {
		log.Println("插入城市设施异常", err)
		return common.New(constant.DBError, "数据库错误")
	}
	err = session.Commit()
	if err != nil {
		log.Println("提交事务异常", err)
		return common.New(constant.DBError, "数据库错误")
	}
	return nil
}

func (c *cityFacilityService) GetByRId(rid int) ([]*data.CityFacility, error) {
	cf := make([]*data.CityFacility, 0)
	err := db.Engine.Table(new(data.CityFacility)).Where("rid=?", rid).Find(&cf)
	if err != nil {
		log.Println(err)
		return cf, common.New(constant.DBError, "数据库错误")
	}
	return cf, nil
}

func (c *cityFacilityService) GetYield(rid int) data.Yield {
	cfs, err := c.GetByRId(rid)
	var y data.Yield
	if err != nil {
		for _, cf := range cfs {
			for _, f := range cf.Facility() {
				if f.GetLevel() > 0 {
					values := gameConfig.FacilityConf.GetValues(f.Type, f.GetLevel())
					additions := gameConfig.FacilityConf.GetAdditions(f.Type)
					for i, aType := range additions {
						if aType == gameConfig.TypeWood {
							y.Wood += values[i]
						} else if aType == gameConfig.TypeGrain {
							y.Grain += values[i]
						} else if aType == gameConfig.TypeIron {
							y.Iron += values[i]
						} else if aType == gameConfig.TypeStone {
							y.Stone += values[i]
						} else if aType == gameConfig.TypeTax {
							y.Gold += values[i]
						}
					}
				}
			}
		}
		log.Println("cityFacilityService GetYield err", err)
	}
	return y
}

func (c *cityFacilityService) GetFacilities(rid, cid int) []data.Facility {
	cf := &data.CityFacility{}
	ok, err := db.Engine.Table(new(data.CityFacility)).Where("rid=? and cityId=?", rid, cid).Get(cf)
	if err != nil {
		log.Println(err)
		return nil
	}
	if ok {
		return cf.Facility()
	}
	return nil
}

func (c *cityFacilityService) UpFacility(rid, cid int, fType int8) (*data.Facility, error) {
	return nil, nil
}
