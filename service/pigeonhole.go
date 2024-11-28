package service

import (
	"Turing-Go/config"
	"Turing-Go/dao"
	"Turing-Go/models"
	"log"
)

func FindPostPigeonhole() models.PigeonholeRes {
	category, err := dao.GetAllCategory()
	if err != nil {
		log.Println(err)
		return models.PigeonholeRes{}
	}
	allPost, err := dao.GetAllPost()
	postsMap := make(map[string][]models.Post)
	for _, post := range allPost {
		at := post.CreateAt.Format("2006-01")
		postsMap[at] = append(postsMap[at], post)
	}
	return models.PigeonholeRes{
		Categorys:    category,
		Viewer:       config.Cfg.Viewer,
		SystemConfig: config.Cfg.System,
		Lines:        postsMap,
	}
}
