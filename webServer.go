package main

import (
	"Turing-Go/config"
	"Turing-Go/server/web"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[WebServer] ")
}
func main() {
	host := config.File.MustValue("web_server", "host", "127.0.0.1")
	port := config.File.MustValue("web_server", "port", "8088")
	router := gin.Default()

	web.Init(router)

	s := &http.Server{
		Addr:           host + ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
