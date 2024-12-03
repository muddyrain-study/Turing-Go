package net

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Server struct {
	addr       string
	router     *Router
	needSecret bool
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}
func (s *Server) NeedSecret(needSecret bool) {
	s.needSecret = needSecret
}
func (s *Server) Router(router *Router) {
	s.router = router
}

func (s *Server) Start() {
	http.HandleFunc("/", s.wsHandler)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		panic(err)
	}
}

var wsUpgrader = websocket.Upgrader{
	// Allow all origins
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Failed to upgrade to websocket: ", err)
	}
	log.Println("New connection from: ", wsConn.RemoteAddr())

	wsServer := NewWsServer(wsConn, s.needSecret)
	wsServer.Router(s.router)
	wsServer.Start()
	wsServer.Handshake()
}
