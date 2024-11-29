package login

import (
	"Turing-Go/net"
	"Turing-Go/server/login/controller"
)

var Router = net.NewRouter()

func init() {
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}
