package model

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Username string
}

func (client *Client) ReadMessageByRoom(name string, rname string, clientsroom map[string]map[*Client]bool, rooms map[string]string, broadcast chan Message, register chan Message, unregister chan Message) {
	defer func() {
		msg := Message{
			Type:     2,
			Username: name,
			Message:  "left",
			Room:     rname,
		}
		client.Conn.Close()
		delete(clientsroom[rname], client)
		if len(clientsroom[rname]) == 0 {
			delete(clientsroom, rname)
			delete(rooms, rname)
		}
		unregister <- msg

	}()
	for {
		_, res, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("message from readMessage", err.Error())
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
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
