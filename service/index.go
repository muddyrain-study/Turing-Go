package service

import (
	"Turing-Go/config"
	"Turing-Go/dao"
	"Turing-Go/models"
	"html/template"
)

func GetAllIndexInfo(page int, pageSize int) (*models.HomeResponse, error) {

	//页面上涉及到的所有的数据，必须有定义
	categorys, err := dao.GetAllCategory()
	if err != nil {
		return nil, err
	}
	posts, err := dao.GetPostPage(page, pageSize)
	if err != nil {
		return nil, err
	}
	var postMores []models.PostMore
	for _, post := range posts {
		categoryName, err := dao.GetCategoryNameById(post.CategoryId)
		if err != nil {
			return nil, err
		}
		userName, err := dao.GetUserNameById(post.UserId)
		content := template.HTML(post.Content)
		if len(content) > 100 {
			content = content[:100]
		}
		postMore := models.PostMore{
			post.Pid,
			post.Title,
			post.Slug,
			content,
			post.CategoryId,
			categoryName,
			post.UserId,
			userName,
			post.ViewCount,
			post.Type,
			models.DateDay(post.CreateAt),
			models.DateDay(post.UpdateAt),
		}
		postMores = append(postMores, postMore)
	}
	homeData := &models.HomeResponse{
		Viewer:    config.Cfg.Viewer,
		Categorys: categorys,
		Posts:     postMores,
		Total:     1,
		Page:      1,
		Pages:     []int{1},
		PageEnd:   true,
	}

	return homeData, nil
}
