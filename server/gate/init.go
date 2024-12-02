package gate

import (
	"Turing-Go/net"
	"Turing-Go/server/gate/controller"
)

var Router = &net.Router{}

func Init() {
	initRouter()
}

func initRouter() {
	controller.GateHandler.Router(Router)
}
