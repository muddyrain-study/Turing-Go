package service

import (
	"Turing-Go/dao"
	"Turing-Go/models"
	"Turing-Go/utils"
	"errors"
)

func Login(userName, passwd string) (*models.Login, error) {
	passwd = utils.Md5Crypt(passwd, "mm")
	user, err := dao.GetUser(userName, passwd)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	uid := user.Uid
	token, err := utils.Award(&uid)
	if err != nil {
		return nil, errors.New("生成token失败")
	}

	var loginRes = &models.Login{
		Token:    token,
		UserInfo: user,
	}
	return loginRes, nil
}
