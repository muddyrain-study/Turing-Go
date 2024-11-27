package models

type Login struct {
	Token    string      `json:"token"`
	UserInfo interface{} `json:"userInfo"`
}
