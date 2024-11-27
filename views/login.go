package views

import (
	"Turing-Go/common"
	"Turing-Go/config"
	"net/http"
)

func (receiver HTMLApi) Login(w http.ResponseWriter, r *http.Request) {
	login := common.Template.Login

	login.WriteData(w, config.Cfg.Viewer)
}
