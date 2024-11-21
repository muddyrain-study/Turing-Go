package main

import (
	"database/sql"
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

func main() {
	defer DB.Close()

	save("test", "男", "123@qq.com")
}
