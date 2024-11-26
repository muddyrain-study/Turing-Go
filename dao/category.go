package dao

import (
	"Turing-Go/models"
	"log"
)

func GetAllCategory() ([]models.Category, error) {
	rows, err := DB.Query("select * from blog_category")

	if err != nil {
		log.Println("blog_category表查询失败", err)
		return nil, err
	}

	var categorys []models.Category

	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.Cid, &category.Name, &category.CreateAt, &category.UpdateAt)
		if err != nil {
			log.Println("blog_category表数据解析失败", err)
			return nil, err
		}
		categorys = append(categorys, category)
	}

	return categorys, nil
}

func GetCategoryNameById(id int) (string, error) {
	row := DB.QueryRow("select name from blog_category where cid = ?", id)

	if row.Err() != nil {
		log.Println("blog_category表查询失败", row.Err())
		return "", nil
	}
	var categoryName string
	_ = row.Scan(&categoryName)
	return categoryName, nil
}
