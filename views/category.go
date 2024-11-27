package views

import (
	"Turing-Go/common"
	"Turing-Go/service"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (receiver HTMLApi) Category(w http.ResponseWriter, r *http.Request) {
	categoryTemplate := common.Template.Category
	path := r.URL.Path
	cIdStr := strings.TrimPrefix(path, "/c/")
	cId, err := strconv.Atoi(cIdStr)
	if err != nil {
		categoryTemplate.WriteError(w, fmt.Errorf("分类ID错误: %s", cIdStr))
	}
	if err := r.ParseForm(); err != nil {
		log.Println("解析请求参数失败", err)
		categoryTemplate.WriteError(w, errors.New("系统错误，请联系管理员！"))
		return
	}
	pageString := r.Form.Get("page")
	page := 1
	if pageString != "" {
		page, _ = strconv.Atoi(pageString)
	}
	pageSize := 10

	categoryResponse, err := service.GetCategoryBy(cId, page, pageSize)
	if err != nil {
		log.Println("获取分类数据失败", err)
		categoryTemplate.WriteError(w, errors.New("系统错误，请联系管理员！"))
		return
	}
	categoryTemplate.WriteData(w, categoryResponse)
}
