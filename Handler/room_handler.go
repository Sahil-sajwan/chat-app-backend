package handler

import (
	"chatapp/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var rooms = make(map[string]string)
var clientsroom = make(map[string]map[*model.Client]bool)
var broadcast = make(chan model.Message)
var register = make(chan model.Message)
var unregister = make(chan model.Message)

func CreateRoomHandler(c *gin.Context) {

	var room model.Room
	c.Bind(&room)
	rname := room.Rname
	rpass := room.Rpass
	if _, ok := rooms[rname]; ok {
		c.JSON(http.StatusConflict, gin.H{
			"message": "room already exists.",
		})
		return
	}
	rooms[rname] = rpass

	c.JSON(http.StatusOK, gin.H{
		"message": "room created",
	})

}

func JoinRoomAuthHandler(c *gin.Context) {
	rname := c.PostForm("rname")
	rpass := c.PostForm("rpass")

	if _, ok := rooms[rname]; !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "requested room does not exist",
		})
		return
	}

	if rpass != rooms[rname] {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "room password is incorrect",
		})
		return
	}

}

func JoinRoomHandler(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	name := c.Param("name")
	rname := c.Query("rname")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {

		return
	}

	client := &model.Client{
		Conn:     conn,
		Username: name,
	}
	if len(clientsroom[rname]) == 0 {
		clientsroom[rname] = make(map[*model.Client]bool)
	}

	clientsroom[rname][client] = true

	msg := model.Message{
		Type:     1,
		Message:  "joined",
		Username: name,
		Room:     rname,
	}
	register <- msg
	client.ReadMessageByRoom(name, rname, clientsroom, rooms, broadcast, register, unregister)
}

func HandleMessagesByRoom() {
	for {

		select {
		case msg := <-register:
			log.Println("from register")
			for client := range clientsroom[msg.Room] {
				err := client.Conn.WriteJSON(msg)
				if err != nil {
					log.Println("register error", err.Error())
					return
				}

			}

		case msg := <-unregister:

			for client := range clientsroom[msg.Room] {
				err := client.Conn.WriteJSON(msg)
				if err != nil {
					log.Println("unregister error", err.Error())
					return
				}

			}

		case msg := <-broadcast:
			for client := range clientsroom[msg.Room] {
				err := client.Conn.WriteJSON(msg)
				if err != nil {
					log.Println("message send error", err.Error())
					return
				}

			}
		}
	}
}
