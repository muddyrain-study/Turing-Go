package main

import (
	"Turing-Go/config"
	"Turing-Go/net"
	"Turing-Go/server/login"
)

func main() {
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	port := config.File.MustValue("login_server", "port", "8004")

	s := net.NewServer(host + ":" + port)

	s.Router(login.Router)
	s.Start()
}
