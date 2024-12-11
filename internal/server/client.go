package server

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Client struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *Client) handleWSConnection() {
	fmt.Println("handling connection")
}
