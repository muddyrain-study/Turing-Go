package views

import (
	"Turing-Go/common"
	"Turing-Go/service"
	"net/http"
)

func (receiver HTMLApi) Pigeonhole(w http.ResponseWriter, r *http.Request) {
	pigeonholeTemplate := common.Template.Pigeonhole

	pigeonholeRes := service.FindPostPigeonhole()

	pigeonholeTemplate.WriteData(w, pigeonholeRes)

}
