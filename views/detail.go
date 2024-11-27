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

func (receiver HTMLApi) Detail(w http.ResponseWriter, r *http.Request) {
	detailTemplate := common.Template.Detail
	path := r.URL.Path
	pIdStr := strings.TrimPrefix(path, "/p/")
	pIdStr = strings.TrimSuffix(pIdStr, ".html")
	pId, err := strconv.Atoi(pIdStr)
	if err != nil {
		detailTemplate.WriteError(w, fmt.Errorf("文章ID错误: %s", pIdStr))
	}
	if err := r.ParseForm(); err != nil {
		log.Println("解析请求参数失败", err)
		detailTemplate.WriteError(w, errors.New("系统错误，请联系管理员！"))
		return
	}
	postRes, err := service.GetPostDetailById(pId)
	if err != nil {
		log.Println("获取文章数据失败", err)
		detailTemplate.WriteError(w, errors.New("系统错误，请联系管理员！"))
		return
	}
	detailTemplate.WriteData(w, postRes)
}
