package controller

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/net"
	"Turing-Go/server/login/model"
	"Turing-Go/server/login/proto"
	"Turing-Go/server/models"
	"Turing-Go/utils"
	"github.com/mitchellh/mapstructure"
	"log"
	"time"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (a *Account) Router(r *net.Router) {
	g := r.Group("account")
	g.AddRouter("login", a.login)
}
func (a *Account) login(req *net.WsMsgReq, resp *net.WsMsgResp) {
	loginReq := &proto.LoginReq{}
	loginResp := &proto.LoginResp{}
	err := mapstructure.Decode(req.Body.Msg, loginReq)
	if err != nil {
		log.Println("login mapstructure.Decode error:", err)
	}
	var user models.User
	ok, err := db.Engine.Table(user).Where("username=?", loginReq.Username).Get(&user)

	if err != nil {
		log.Println("login db.Engine.Table(user).Where error:", err)
		return
	}
	if !ok {
		resp.Body.Code = constant.UserNotExist
		return
	}
	pwd := utils.Password(loginReq.Password, user.Passcode)
	if pwd != user.Passwd {
		resp.Body.Code = constant.PwdIncorrect
		return
	}
	token, err := utils.Award(user.UId)
	if err != nil {
		log.Println("login utils.Award error:", err)
	}
	resp.Body.Code = constant.OK
	loginResp.UId = user.UId
	loginResp.Username = user.Username
	loginResp.Session = token
	loginResp.Password = ""
	resp.Body.Msg = loginResp

	ul := &model.LoginHistory{
		UId:      user.UId,
		Ip:       loginReq.Ip,
		CTime:    time.Now(),
		Hardware: loginReq.Hardware,
		State:    model.Login,
	}
	_, err = db.Engine.Table(ul).Insert(ul)
	if err != nil {
		log.Println("login db.Engine.Table(ul).Insert error:", err)
	}

	ll := &model.LoginLast{
		UId: user.UId,
	}
	ok, err = db.Engine.Table(ll).Where("uid = ?", user.UId).Get(ll)
	if err != nil {
		log.Println("login db.Engine.Table(ll).Where error:", err)
	}
	if ok {
		ll.IsLogout = 0
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		ll.Hardware = loginReq.Hardware
		db.Engine.Table(ll).Where("uid = ?", user.UId).Update(ll)
	} else {
		ll.IsLogout = 0
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		ll.Hardware = loginReq.Hardware
		ll.UId = user.UId
		_, err := db.Engine.Table(ll).Insert(ll)
		if err != nil {
			log.Println("login db.Engine.Table(ll).Insert error:", err)
		}
	}
	//缓存一下 此用户和当前的ws连接
	net.Mgr.UserLogin(req.Conn, user.UId, token)
}
