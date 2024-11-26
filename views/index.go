package views

import (
	"Turing-Go/common"
	"Turing-Go/config"
	"Turing-Go/models"
	"net/http"
)

func (receiver HTMLApi) Index(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/html")
	var homeData = new(models.HomeResponse)
	//页面上涉及到的所有的数据，必须有定义
	var categorys = []models.Category{
		{
			Cid:  1,
			Name: "go",
		},
	}
	var posts = []models.PostMore{
		{
			Pid:          1,
			Title:        "go博客",
			Content:      "内容",
			UserName:     "张三",
			ViewCount:    123,
			CreateAt:     "2022-02-20",
			CategoryId:   1,
			CategoryName: "go",
			Type:         0,
		},
	}
	homeData = &models.HomeResponse{
		Viewer:    config.Cfg.Viewer,
		Categorys: categorys,
		Posts:     posts,
		Total:     1,
		Page:      1,
		Pages:     []int{1},
		PageEnd:   true,
	}
	index := common.Template.Index
	index.WriteData(w, homeData)
}
