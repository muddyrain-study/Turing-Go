package controller

import (
	"Turing-Go/constant"
	"Turing-Go/net"
	"Turing-Go/server/common"
	"Turing-Go/server/game/logic"
	"Turing-Go/server/game/model"
	"Turing-Go/server/game/model/data"
)

var WarController = &warController{}

type warController struct {
}

func (w *warController) InitRouter(router *net.Router) {
	r := router.Group("war")
	r.AddRouter("report", w.report)
}

func (w *warController) report(req *net.WsMsgReq, resp *net.WsMsgResp) {
	respObj := &model.WarReportRsp{}
	resp.Body.Seq = req.Body.Seq
	resp.Body.Name = req.Body.Name
	role, _ := req.Conn.GetProperty("role")
	r := role.(data.Role)
	reports, err := logic.WarService.GetWarReports(r.RId)
	if err != nil {
		resp.Body.Code = err.(*common.MyError).Code()
		return
	}
	respObj.List = reports
	resp.Body.Code = constant.OK
	resp.Body.Msg = respObj
}
