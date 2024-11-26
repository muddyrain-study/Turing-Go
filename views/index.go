package views

import (
	"Turing-Go/common"
	"Turing-Go/service"
	"errors"
	"log"
	"net/http"
	"strconv"
)

func (receiver HTMLApi) Index(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/html")
	index := common.Template.Index
	if err := r.ParseForm(); err != nil {
		log.Println("解析请求参数失败", err)
		index.WriteError(w, errors.New("系统错误，请联系管理员！"))
		return
	}
	pageString := r.Form.Get("page")
	page := 1
	if pageString != "" {
		page, _ = strconv.Atoi(pageString)
	}
	pageSize := 10
	homeData, err := service.GetAllIndexInfo(page, pageSize)

	if err != nil {
		log.Println("获取首页数据失败", err)
		index.WriteError(w, errors.New("系统错误，请联系管理员！"))
	}
	index.WriteData(w, homeData)
}
