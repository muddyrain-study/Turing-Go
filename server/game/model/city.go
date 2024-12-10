package model

type FacilitiesReq struct {
	CityId int `json:"cityId"`
}

type FacilitiesRsp struct {
	CityId     int        `json:"cityId"`
	Facilities []Facility `json:"facilities"`
}
type Facility struct {
	Name   string `json:"name"`
	Level  int8   `json:"level"`
	Type   int8   `json:"type"`
	UpTime int64  `json:"up_time"` //升级的时间戳，0表示该等级已经升级完成了
}
