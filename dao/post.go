package dao

import (
	"Turing-Go/models"
	"log"
)

func GetPostPage(page, pageSize int) ([]models.Post, error) {
	page = (page - 1) * pageSize
	rows, err := DB.Query("select * from blog_post limit ?,?", page, pageSize)

	if err != nil {
		log.Println("blog_post表查询失败", err)
		return nil, err
	}

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.Pid,
			&post.Title,
			&post.Content,
			&post.Markdown,
			&post.CategoryId,
			&post.UserId,
			&post.ViewCount,
			&post.Type,
			&post.Slug,
			&post.CreateAt,
			&post.UpdateAt,
		)
		if err != nil {
			log.Println("blog_post表数据解析失败", err)
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func CountGetAllPost() int {
	rows := DB.QueryRow("select count(1) from blog_post")
	var count int
	_ = rows.Scan(&count)
	return count
}

func CountGetAllPostByCategoryId(cId int) int {
	rows := DB.QueryRow("select count(1) from blog_post where category_id = ?", cId)
	var count int
	_ = rows.Scan(&count)
	return count
}

func GetPostPageByCategoryId(page, pageSize, cId int) ([]models.Post, error) {
	page = (page - 1) * pageSize
	rows, err := DB.Query("select * from blog_post where category_id = ? limit ?,?", cId, page, pageSize)

	if err != nil {
		log.Println("blog_post表查询失败", err)
		return nil, err
	}

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.Pid,
			&post.Title,
			&post.Content,
			&post.Markdown,
			&post.CategoryId,
			&post.UserId,
			&post.ViewCount,
			&post.Type,
			&post.Slug,
			&post.CreateAt,
			&post.UpdateAt,
		)
		if err != nil {
			log.Println("blog_post表数据解析失败", err)
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
