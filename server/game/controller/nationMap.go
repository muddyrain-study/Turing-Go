package controller

import (
	"Turing-Go/constant"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/logic"
	"Turing-Go/server/game/middleware"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
	"github.com/mitchellh/mapstructure"
)

var NationMapController = &nationMapController{}

type nationMapController struct {
}

func (r *nationMapController) InitRouter(router *net.Router) {
	g := router.Group("nationMap")
	g.AddRouter("config", r.config)
	g.AddRouter("scanBlock", r.scanBlock, middleware.CheckRole())
}

func (r *nationMapController) config(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	rspObj := &model.ConfigRsp{}

	m := gameConfig.MapBuildConf.Cfg
	rspObj.Confs = make([]model.Conf, len(m))
	for index, v := range m {
		rspObj.Confs[index].Type = v.Type
		rspObj.Confs[index].Name = v.Name
		rspObj.Confs[index].Level = v.Level
		rspObj.Confs[index].Defender = v.Defender
		rspObj.Confs[index].Durable = v.Durable
		rspObj.Confs[index].Grain = v.Grain
		rspObj.Confs[index].Iron = v.Iron
		rspObj.Confs[index].Stone = v.Stone
		rspObj.Confs[index].Wood = v.Wood
	}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK

}

func (*nationMapController) scanBlock(req *net.WsMsgReq, resp *net.WsMsgResp) {
	reqObj := &model.ScanBlockReq{}
	respObj := &model.ScanRsp{}

	err := mapstructure.Decode(req.Body.Msg, reqObj)
	if err != nil {
		resp.Body.Code = constant.InvalidParam
		return
	}
	resp.Body.Name = req.Body.Name
	resp.Body.Seq = req.Body.Seq
	resp.Body.Code = constant.OK

	mrb, err := logic.RoleBuildService.ScanBlock(reqObj)
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	respObj.MRBuilds = mrb
	mrc, err := logic.RoleCityService.ScanBlock(reqObj)
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	respObj.MCBuilds = mrc
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		resp.Body.Code = constant.SessionInvalid
		return
	}
	rid := role.(data.Role).RId
	armies, err := logic.ArmyService.ScanBlock(rid, reqObj)
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	respObj.Armys = armies
	resp.Body.Msg = respObj
}
