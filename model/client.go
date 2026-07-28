package model

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Username string
}

func (client *Client) ReadMessageByRoom(name string, rname string, broadcast chan Message, unregister chan Message, cleanup func()) {
	defer func() {
		msg := Message{
			Type:     2,
			Username: name,
			Message:  "left",
			Room:     rname,
		}
		client.Conn.Close()
		cleanup()
		unregister <- msg

	}()
	for {
		_, res, err := client.Conn.ReadMessage()
		if err != nil {

			return

		}
		msg := Message{
			Type:     3,
			Username: name,
			Message:  string(res),
			Room:     rname,
		}
		broadcast <- msg

	}
}
