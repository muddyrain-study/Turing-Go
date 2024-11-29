package controller

import (
	"Turing-Go/net"
	"Turing-Go/server/login/proto"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (a *Account) Router(r *net.Router) {
	g := r.Group("account")
	g.AddRouter("login", a.login)
}
func (a *Account) login(req *net.WsMsgReq, resp *net.WsMsgResp) {
	resp.Body.Code = 0
	loginResp := &proto.LoginResp{}
	loginResp.UId = 1
	loginResp.Username = "admin"
	loginResp.Session = "AS"
	loginResp.Password = ""
	resp.Body.Msg = loginResp
}
