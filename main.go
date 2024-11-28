package main

import (
	"Turing-Go/config"
	"fmt"
)

func init() {
}
func main() {
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	fmt.Println(host)
}
