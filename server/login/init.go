package login

import (
	"Turing-Go/db"
	"Turing-Go/net"
	"Turing-Go/server/login/controller"
)

var Router = net.NewRouter()

func Init() {
	db.TestDB()
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}
