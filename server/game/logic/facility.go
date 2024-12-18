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
	"time"
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

func (c *cityFacilityService) GetFacilities1(rid, cid int) []*data.Facility {
	cf := &data.CityFacility{}
	ok, err := db.Engine.Table(new(data.CityFacility)).Where("rid=? and cityId=?", rid, cid).Get(cf)
	if err != nil {
		log.Println(err)
		return nil
	}
	if ok {
		return cf.Facility1()
	}
	return nil
}

func (c *cityFacilityService) UpFacility(rid, cid int, fType int8) (*data.Facility, error) {
	facs := c.GetFacilities1(rid, cid)
	result := &data.Facility{}
	for _, fac := range facs {
		if fac.Type == fType {
			// 判断是否可以升级 该设施是否在升级中，是否自愿够
			if !fac.CanUp() {
				return nil, common.New(constant.UpError, "升级失败")
			}
			maxLevel := fac.GetMaxLevel(fType)
			if fac.GetLevel() >= int8(maxLevel) {
				return nil, common.New(constant.UpError, "已经是最高等级")
			}
			// 先判断资源使用多少，再判断是否够
			need := gameConfig.FacilityConf.Need(fType, fac.GetLevel()+1)
			ok := RoleResService.TryUseNeed(rid, need)
			if !ok {
				return nil, common.New(constant.UpError, "资源不足")
			}
			fac.UpTime = time.Now().Unix()
			result = fac
		}
	}
	jsonByte, _ := json.Marshal(facs)
	cfac := c.Get(rid, cid)
	cfac.Facilities = string(jsonByte)
	cfac.SyncExecute()
	return result, nil
}

func (c *cityFacilityService) Get(rid int, cid int) *data.CityFacility {
	cf := &data.CityFacility{}
	ok, err := db.Engine.Table(new(data.CityFacility)).Where("rid=? and cityId=?", rid, cid).Get(cf)
	if err != nil {
		log.Println(err)
		return nil
	}
	if ok {
		return cf
	}
	return nil
}
func (c *cityFacilityService) GetByCid(cid int) *data.CityFacility {
	cf := &data.CityFacility{}
	ok, err := db.Engine.Table(new(data.CityFacility)).Where("cityId=?", cid).Get(cf)
	if err != nil {
		log.Println(err)
		return nil
	}
	if ok {
		return cf
	}
	return nil
}

func (c *cityFacilityService) GetFacilityLevel(cid int, fType int8) int8 {
	cf := c.GetByCid(cid)
	if cf == nil {
		return 0
	}
	facs := cf.Facility1()
	for _, v := range facs {
		if v.Type == fType {
			return v.GetLevel()
		}
	}
	return 0
}

func (c *cityFacilityService) GetCost(cid int) int8 {
	cf := c.GetByCid(cid)
	facility := cf.Facility()
	var cost int
	for _, f := range facility {
		if f.GetLevel() > 0 {
			values := gameConfig.FacilityConf.GetValues(f.Type, f.GetLevel())
			additions := gameConfig.FacilityConf.GetAdditions(f.Type)
			for i, aType := range additions {
				if aType == gameConfig.TypeCost {
					cost += values[i]
				}
			}
		}
	}
	return int8(cost)
}

func (c *cityFacilityService) GetCapacity(rid int) int {
	cfs, err := c.GetByRId(rid)
	var cap int
	if err != nil {
		for _, cf := range cfs {
			for _, f := range cf.Facility() {
				if f.GetLevel() > 0 {
					values := gameConfig.FacilityConf.GetValues(f.Type, f.GetLevel())
					additions := gameConfig.FacilityConf.GetAdditions(f.Type)
					for i, aType := range additions {
						if aType == gameConfig.TypeWarehouseLimit {
							cap += values[i]
						}
					}
				}
			}
		}
		log.Println("cityFacilityService GetYield err", err)
	}
	return cap
}
