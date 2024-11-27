package dao

import (
	"Turing-Go/models"
	"log"
)

func GetUserNameById(id int) (string, error) {
	row := DB.QueryRow("select userName from blog_user where uid = ?", id)

	if row.Err() != nil {
		log.Println("blog_user表查询失败", row.Err())
		return "", nil
	}
	var userName string
	_ = row.Scan(&userName)
	return userName, nil
}

func GetUser(userName, passwd string) (*models.User, error) {
	row := DB.QueryRow("select * from blog_user where user_name = ? and passwd = ? limit 1", userName, passwd)

	if row.Err() != nil {
		log.Println("blog_user表查询失败", row.Err())
		return nil, nil
	}
	var user = &models.User{}
	err := row.Scan(&user.Uid, &user.UserName, &user.Password, &user.Avatar, &user.CreateAt, &user.UpdateAt)
	if err != nil {
		log.Println("解析用户数据失败", err)
		return nil, err
	}
	return user, nil
}
