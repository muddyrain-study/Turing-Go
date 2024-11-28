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

func GetAllPost() ([]models.Post, error) {
	rows, err := DB.Query("select * from blog_post")
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

func GetPostById(pid int) (*models.Post, error) {
	row := DB.QueryRow("select * from blog_post where pid = ?", pid)
	if row.Err() != nil {
		log.Println("blog_post表查询失败", row.Err())
		return nil, nil
	}
	var post models.Post
	err := row.Scan(
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
	return &post, nil
}

func SavePost(post *models.Post) {
	log.Println("post", post)
	res, err := DB.Exec(
		"insert into blog_post(title, content, markdown, category_id, user_id, view_count, type, slug, createAt, updateAt)"+
			" values(?,?,?,?,?,?,?,?,?,?)",
		post.Title,
		post.Content,
		post.Markdown,
		post.CategoryId,
		post.UserId,
		post.ViewCount,
		post.Type,
		post.Slug,
		post.CreateAt.Format("2006-01-02 15:04:05"),
		post.UpdateAt.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		log.Println("blog_post表插入失败", err)
		return
	}
	pid, _ := res.LastInsertId()
	log.Println("pid", pid)
	post.Pid = int(pid)
}

func UpdatePost(post *models.Post) {
	_, err := DB.Exec(
		"update blog_post set title = ?, content = ?, markdown = ?, category_id = ?, user_id = ?, type = ?, slug = ?, updateAt = ? where pid = ?",
		post.Title,
		post.Content,
		post.Markdown,
		post.CategoryId,
		post.UserId,
		post.Type,
		post.Slug,
		post.UpdateAt.Format("2006-01-02 15:04:05"),
		post.Pid,
	)
	if err != nil {
		log.Println("blog_post表更新失败", err)
		return
	}
}
