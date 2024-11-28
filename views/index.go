package views

import (
	"Turing-Go/common"
	"Turing-Go/service"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (receiver HTMLApi) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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

	path := r.URL.Path
	slug := strings.TrimPrefix(path, "/")
	homeData, err := service.GetAllIndexInfo(slug, page, pageSize)

	if err != nil {
		log.Println("获取首页数据失败", err)
		index.WriteError(w, errors.New("系统错误，请联系管理员！"))
	}
	index.WriteData(w, homeData)
}
