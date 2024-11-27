package service

import (
	"Turing-Go/config"
	"Turing-Go/dao"
	"Turing-Go/models"
	"fmt"
	"html/template"
)

func GetPostDetailById(pid int) (*models.PostRes, error) {
	post, err := dao.GetPostById(pid)
	fmt.Println(post, err)
	if err != nil {
		return &models.PostRes{}, err
	}
	categoryName, err := dao.GetCategoryNameById(post.CategoryId)
	if err != nil {
		return nil, err
	}
	userName, err := dao.GetUserNameById(post.UserId)
	postMore := models.PostMore{
		Pid:          post.Pid,
		Title:        post.Title,
		Slug:         post.Slug,
		Content:      template.HTML(post.Content),
		CategoryId:   post.CategoryId,
		CategoryName: categoryName,
		UserId:       post.UserId,
		UserName:     userName,
		ViewCount:    post.ViewCount,
		Type:         post.Type,
		CreateAt:     models.DateDay(post.CreateAt),
		UpdateAt:     models.DateDay(post.UpdateAt),
	}
	var postRes = &models.PostRes{
		Viewer:       config.Cfg.Viewer,
		SystemConfig: config.Cfg.System,
		Article:      postMore,
	}

	return postRes, nil
}
