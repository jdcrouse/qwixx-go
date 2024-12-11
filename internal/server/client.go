package server

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func (c *Client) handleWSConnection() {
	fmt.Println("handling connection")
}
