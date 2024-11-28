package router

import (
	"Turing-Go/api"
	"Turing-Go/views"
	"net/http"
)

func Router() {
	http.HandleFunc("/", views.HTML.Index)
	http.HandleFunc("/c/", views.HTML.Category)
	http.HandleFunc("/p/", views.HTML.Detail)
	http.HandleFunc("/login", views.HTML.Login)
	http.HandleFunc("/writing", views.HTML.Writing)
	http.HandleFunc("/pigeonhole", views.HTML.Pigeonhole)
	http.HandleFunc("/api/v1/post", api.API.SaveAndUpdatePost)
	http.HandleFunc("/api/v1/post/search", api.API.SearchPost)
	http.HandleFunc("/api/v1/post/", api.API.GetPost)
	http.HandleFunc("/api/v1/qiniu/token", api.API.QiNiuToken)
	http.HandleFunc("/api/v1/login", api.API.Login)
	http.Handle("/resource/", http.StripPrefix("/resource/", http.FileServer(http.Dir("public/resource/"))))
}
