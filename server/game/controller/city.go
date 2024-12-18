package controller

import (
	"Turing-Go/constant"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/logic"
	"Turing-Go/server/game/middleware"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"github.com/mitchellh/mapstructure"
)

var CityController = &cityController{}

type cityController struct {
}

func (g *cityController) InitRouter(router *net.Router) {
	r := router.Group("city")
	r.Use(middleware.Log())
	r.AddRouter("facilities", g.facilities, middleware.CheckRole())
	r.AddRouter("upFacility", g.upFacility, middleware.CheckRole())
}

func (g *cityController) facilities(req *net.WsMsgReq, resp *net.WsMsgResp) {
	reqObj := &model.FacilitiesReq{}
	respObj := &model.FacilitiesRsp{}
	err := mapstructure.Decode(req.Body.Msg, reqObj)
	if err != nil {
		resp.Body.Code = constant.InvalidParam
		return
	}
	resp.Body.Msg = respObj
	respObj.CityId = reqObj.CityId
	resp.Body.Code = constant.OK

	r, _ := req.Conn.GetProperty("role")
	city, ok := logic.RoleCityService.Get(reqObj.CityId)
	if ok == false {
		resp.Body.Code = constant.CityNotExist
		return
	}

	role := r.(data.Role)
	if city.RId != role.RId {
		resp.Body.Code = constant.CityNotMe
		return
	}
	fac := logic.CityFacilityService.GetFacilities(role.RId, reqObj.CityId)
	if fac == nil {
		resp.Body.Code = constant.DBError
		return
	}
	respObj.CityId = reqObj.CityId
	respObj.Facilities = make([]model.Facility, len(fac))
	for i, v := range fac {
		respObj.Facilities[i].Name = v.Name
		respObj.Facilities[i].Level = v.GetLevel()
		respObj.Facilities[i].Type = v.Type
		respObj.Facilities[i].UpTime = v.UpTime
	}
}

func (g *cityController) upFacility(req *net.WsMsgReq, resp *net.WsMsgResp) {
	reqObj := &model.UpFacilityReq{}
	respObj := &model.UpFacilityRsp{}
	err := mapstructure.Decode(req.Body.Msg, reqObj)
	if err != nil {
		resp.Body.Code = constant.InvalidParam
		return
	}
	resp.Body.Msg = respObj
	respObj.CityId = reqObj.CityId
	resp.Body.Code = constant.OK

	r, _ := req.Conn.GetProperty("role")
	city, ok := logic.RoleCityService.Get(reqObj.CityId)
	if ok == false {
		resp.Body.Code = constant.CityNotExist
		return
	}

	role := r.(data.Role)
	if city.RId != role.RId {
		resp.Body.Code = constant.CityNotMe
		return
	}

	facs := logic.CityFacilityService.GetFacilities(role.RId, reqObj.CityId)
	if facs == nil {
		resp.Body.Code = constant.CityNotExist
		return
	}
	fac, err := logic.CityFacilityService.UpFacility(role.RId, reqObj.CityId, reqObj.FType)

	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	respObj.Facility.Name = fac.Name
	respObj.Facility.Name = fac.Name
	respObj.Facility.Level = fac.GetLevel()
	respObj.Facility.Type = fac.Type
	respObj.Facility.UpTime = fac.UpTime
	if roleRes := logic.RoleResService.GetRoleRes(role.RId); roleRes != nil {
		respObj.RoleRes = roleRes.ToModel().(model.RoleRes)
	}
}
