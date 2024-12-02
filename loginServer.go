package main

import (
	"Turing-Go/config"
	"Turing-Go/net"
	"Turing-Go/server/login"
	"log"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[LoginServer] ")
}

func main() {
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	port := config.File.MustValue("login_server", "port", "8004")

	s := net.NewServer(host + ":" + port)
	login.Init()
	s.Router(login.Router)
	s.Start()
}
