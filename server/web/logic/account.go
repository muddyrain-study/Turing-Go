package logic

import (
	"Turing-Go/constant"
	"Turing-Go/db"
	"Turing-Go/server/common"
	"Turing-Go/server/models"
	"Turing-Go/server/web/model"
	"Turing-Go/utils"
	"log"
	"time"
)

var DefaultAccountLogic = &AccountLogic{}

type AccountLogic struct{}

func (a *AccountLogic) Register(req *model.RegisterReq) error {
	username := req.Username
	user := models.User{}

	ok, err := db.Engine.Table(user).Where("username = ?", username).Get(&user)
	if err != nil {
		log.Println("查询用户异常:", err)
		return common.New(constant.DBError, "数据库查询异常")
	}
	if ok {
		return common.New(constant.UserExist, "用户已存在")
	} else {
		user.Username = username
		user.Hardware = req.Hardware
		user.Ctime = time.Now()
		user.Mtime = time.Now()
		user.Passcode = utils.RandSeq(6)
		user.Passwd = utils.Password(req.Password, user.Passcode)

		_, err := db.Engine.Table(user).Insert(&user)
		if err != nil {
			log.Println("插入用户异常:", err)
			return common.New(constant.DBError, "数据库插入异常")
		}
		return nil
	}
}
