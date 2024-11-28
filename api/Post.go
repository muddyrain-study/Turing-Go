package api

import (
	"Turing-Go/common"
	"Turing-Go/dao"
	"Turing-Go/models"
	"Turing-Go/service"
	"Turing-Go/utils"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (receiver Api) GetPost(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	pIdStr := strings.TrimPrefix(url, "/api/v1/post/")
	pId, err := strconv.Atoi(pIdStr)
	if err != nil {
		common.Error(w, errors.New("id is invalid"))
		return
	}
	post, err := dao.GetPostById(pId)
	if err != nil {
		common.Error(w, errors.New("post is not found"))
		return
	}
	common.Success(w, post)
}

func (receiver Api) SaveAndUpdatePost(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	_, claims, err := utils.ParseToken(token)
	if err != nil {
		common.Error(w, errors.New("token is invalid"))
		return
	}
	method := r.Method

	switch method {
	case http.MethodPost:
		params := common.GetRequestJsonParams(r)
		cId := params["categoryId"].(string)
		categoryId, _ := strconv.Atoi(cId)
		content := params["content"].(string)
		markdown := params["markdown"].(string)
		slug := params["slug"].(string)
		title := params["title"].(string)
		var postType int
		if val, ok := params["type"]; ok {
			if floatVal, ok := val.(float64); ok {
				postType = int(floatVal)
			}
		} else {
			postType = 0
		}
		post := &models.Post{
			Pid:        -1,
			Title:      title,
			Slug:       slug,
			Content:    content,
			Markdown:   markdown,
			CategoryId: categoryId,
			ViewCount:  0,
			Type:       postType,
			UserId:     claims.Uid,
			CreateAt:   time.Now(),
			UpdateAt:   time.Now(),
		}
		service.SavePost(post)
		common.Success(w, post)
	case http.MethodPut:
		// update
		params := common.GetRequestJsonParams(r)
		categoryId := int(params["categoryId"].(float64))
		content := params["content"].(string)
		markdown := params["markdown"].(string)
		slug := params["slug"].(string)
		title := params["title"].(string)
		pId := int(params["pid"].(float64))
		postType := int(params["type"].(float64))
		post := &models.Post{
			Pid:        pId,
			Title:      title,
			Slug:       slug,
			Content:    content,
			Markdown:   markdown,
			CategoryId: categoryId,
			Type:       postType,
			UserId:     claims.Uid,
			UpdateAt:   time.Now(),
		}
		service.UpdatePost(post)
		common.Success(w, post)
	}
}
