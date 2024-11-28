package views

import (
	"Turing-Go/common"
	"Turing-Go/service"
	"net/http"
)

func (receiver HTMLApi) Writing(w http.ResponseWriter, r *http.Request) {
	writingTemplate := common.Template.Writing
	wr := service.Writing()
	writingTemplate.WriteData(w, wr)
}
