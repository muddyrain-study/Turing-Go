package api

import (
	"Turing-Go/common"
	"Turing-Go/service"
	"net/http"
)

func (receiver Api) Login(w http.ResponseWriter, r *http.Request) {
	params := common.GetRequestJsonParams(r)
	username := params["username"].(string)
	passwd := params["passwd"].(string)
	login, err := service.Login(username, passwd)
	if err != nil {
		common.Error(w, err)
	}
	common.Success(w, login)
}
