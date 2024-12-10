package data

import (
	"Turing-Go/server/game/gameConfig"
	"encoding/json"
	"time"
)

type Facility struct {
	Name         string `json:"name"`
	PrivateLevel int8   `json:"level"` //等级，外部读的时候不能直接读，要用GetLevel
	Type         int8   `json:"type"`
	UpTime       int64  `json:"up_time"` //升级的时间戳，0表示该等级已经升级完成了
}

func (f Facility) GetLevel() int8 {
	if f.UpTime > 0 {
		cur := time.Now().Unix()
		cost := gameConfig.FacilityConf.CostTime(f.Type, f.PrivateLevel+1)
		if cur >= f.UpTime+int64(cost) {
			f.PrivateLevel += 1
			f.UpTime = 0
		}
	}
	return f.PrivateLevel
}

type CityFacility struct {
	Id         int    `xorm:"id pk autoincr"`
	RId        int    `xorm:"rid"`
	CityId     int    `xorm:"cityId"`
	Facilities string `xorm:"facilities"`
}

func (c *CityFacility) TableName() string {
	return "city_facility"
}

func (c *CityFacility) Facility() []Facility {
	facilities := make([]Facility, 0)
	err := json.Unmarshal([]byte(c.Facilities), &facilities)
	if err != nil {
		return nil
	}
	return facilities
}
