package dao

import "log"

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
