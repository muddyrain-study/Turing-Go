package controller

import (
	"Turing-Go/constant"
	"Turing-Go/net"
	"Turing-Go/server/game/gameConfig"
	"Turing-Go/server/game/model"
)

var NationMapController = &nationMapController{}

type nationMapController struct {
}

func (r *nationMapController) InitRouter(router *net.Router) {
	g := router.Group("nationMap")
	g.AddRouter("config", r.config)
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
