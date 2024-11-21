package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var DB *sql.DB

func init() {
	// 执行 main 之前，会先执行 init
	db, err := sql.Open("mysql", "root:mm001030@tcp(127.0.0.1:3306)/go_test")
	if err != nil {
		log.Println("open database error:", err)
		panic(err)
	}
	// 设置最大连接数
	db.SetMaxIdleConns(5)
	// 设置最大打开连接数
	db.SetMaxOpenConns(100)
	// 设置连接的最大存活时间
	db.SetConnMaxLifetime(time.Minute * 3)
	// 设置连接的最大空闲时间
	db.SetConnMaxIdleTime(time.Minute * 1)

	err = db.Ping()
	if err != nil {
		log.Println("open database error:", err)
		db.Close()
		panic(err)
	}
	DB = db

	log.Println("open database success")
}

func save(username string, sex string, email string) {
	res, err := DB.Exec("insert into user (username,sex,email) values (?,?,?)", username, sex, email)
	if err != nil {
		log.Println("insert data error:", err)
		return
	}
	id, _ := res.LastInsertId()
	log.Println("insert success, id:", id)
}

type User struct {
	UserId   int    `db:"user_id"`
	Username string `db:"username"`
	Sex      string `db:"sex"`
	Email    string `db:"email"`
}

func query(id int) (*User, error) {
	rows, err := DB.Query("select * from user where user_id = ? limit 1", id)
	if err != nil {
		log.Println("query data error:", err)
		return nil, err
	}
	user := new(User)
	for rows.Next() {
		err = rows.Scan(&user.UserId, &user.Username, &user.Sex, &user.Email)
		if err != nil {
			log.Println("scan data error:", err)
			return nil, err
		}

	}
	return user, nil
}

func update(id int, username string) {
	res, err := DB.Exec("update user set username = ? where user_id = ?", username, id)
	if err != nil {
		log.Println("update data error:", err)
		return
	}
	num, _ := res.RowsAffected()
	log.Println("update success, num:", num)
}
func insertTx(username string) {
	tx, err := DB.Begin()
	if err != nil {
		log.Println("begin tx error:", err)
		return
	}
	res, err := tx.Exec("insert into user (username,sex,email) values (?,?,?)", username, "男", "456@qq.com")
	if err != nil {
		log.Println("insert data error:", err)
		tx.Rollback()
		return
	}
	if username == "tx" {
		log.Println("rollback")
		tx.Rollback()
		return
	}
	id, _ := res.LastInsertId()
	fmt.Println("insert success, id:", id)
	_ = tx.Commit()
}
func main() {
	defer DB.Close()
	insertTx("tx")
}
