package controller

import (
	"Turing-Go/constant"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/logic"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
)

var SkillController = &skillController{}

type skillController struct {
}

func (w *skillController) InitRouter(router *net.Router) {
	r := router.Group("skill")
	r.AddRouter("list", w.list)
}

func (w *skillController) list(req *net.WsMsgReq, resp *net.WsMsgResp) {
	respObj := &model.SkillListRsp{}
	resp.Body.Msg = respObj
	resp.Body.Code = constant.OK
	resp.Body.Seq = req.Body.Seq
	resp.Body.Name = req.Body.Name

	role, _ := req.Conn.GetProperty("role")
	r := role.(data.Role)
	skills, err := logic.SkillService.GetSkills(r.RId)
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	respObj.List = skills
	resp.Body.Msg = respObj
}
