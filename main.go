package main

import (
	"Turing-Go/common"
	"Turing-Go/router"
	"log"
	"net/http"
)

func init() {
	common.LoadTemplate()
}
func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	router.Router()

	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server Listen Error: %s", err)
	}
}
