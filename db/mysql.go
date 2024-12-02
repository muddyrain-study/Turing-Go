package db

import (
	"Turing-Go/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

var Engine *xorm.Engine

func TestDB() {
	mysqlConfig, err := config.File.GetSection("mysql")
	if err != nil {
		log.Fatal("mysql config error :", err)
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConfig["user"],
		mysqlConfig["password"],
		mysqlConfig["host"],
		mysqlConfig["port"],
		mysqlConfig["dbname"],
		mysqlConfig["charset"],
	)
	Engine, err = xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		log.Fatal("mysql connect error :", err)
	}
	err = Engine.Ping()
	if err != nil {
		log.Fatal("mysql ping error :", err)
	}
	maxIdle := config.File.MustInt("mysql", "max_idle", 2)
	maxConn := config.File.MustInt("mysql", "max_conn", 2)
	Engine.SetMaxIdleConns(maxIdle)
	Engine.SetMaxOpenConns(maxConn)
	Engine.ShowSQL(true)
	log.Println("mysql connect success")
}
