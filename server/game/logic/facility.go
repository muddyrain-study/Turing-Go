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
	ok, err := db.Engine.Table(cf).Where("cid=?").Get(cf)
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
	return nil
}
