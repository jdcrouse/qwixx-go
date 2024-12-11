package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Server interface {
	Start(settings Settings) error
}

type serverImpl struct {
	wsUpgrader websocket.Upgrader
}

func New() Server {
	return &serverImpl{}
}

type Settings struct {
	Endpoint string
}

func (s *serverImpl) Start(settings Settings) error {
	fmt.Println("server starting")
	http.HandleFunc("/ws", s.serveWs)
	return http.ListenAndServe(settings.Endpoint, nil)
}

func (s *serverImpl) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := s.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{conn: conn}
	go client.handleWSConnection()
}
