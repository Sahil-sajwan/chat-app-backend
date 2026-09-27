package handler

import (
	"chatapp/db"
	"chatapp/model"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var rooms = make(map[string]string)
var clientsroom = make(map[string]map[*model.Client]bool)
var roomMu sync.RWMutex
var broadcast = make(chan model.Message)
var register = make(chan model.Message)
var unregister = make(chan model.Message)
var joinTokens = make(map[string]joinToken)
var joinTokensMu sync.Mutex
var disconnectTimers = make(map[string]*time.Timer)
var disconnectMu sync.Mutex

const joinTokenTTL = 5 * time.Minute
const disconnectGracePeriod = 4 * time.Second

type joinToken struct {
	Room      string
	ExpiresAt time.Time
}

func CreateRoomHandler(c *gin.Context) {
	var room model.Room
	c.Bind(&room)
	rname := room.Rname
	rpass := room.Rpass
	if !createRoom(rname, rpass) {
		c.JSON(http.StatusConflict, gin.H{
			"message": "room already exists.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "room created",
	})
}

func JoinRoomAuthHandler(c *gin.Context) {
	rname := c.PostForm("rname")
	rpass := c.PostForm("rpass")

	ok, passwordMatches := validateRoomPassword(rname, rpass)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "requested room does not exist",
		})
		return
	}

	if !passwordMatches {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "room password is incorrect",
		})
		return
	}

	token, err := newJoinToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not create join token",
		})
		return
	}

	joinTokensMu.Lock()
	joinTokens[token] = joinToken{
		Room:      rname,
		ExpiresAt: time.Now().Add(joinTokenTTL),
	}
	joinTokensMu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message":    "room authenticated",
		"token":      token,
		"expires_in": int(joinTokenTTL.Seconds()),
	})
}

func JoinRoomHandler(c *gin.Context) {
	name := c.Param("name")
	rname := c.Query("rname")
	token := c.Query("token")

	if !consumeJoinToken(token, rname) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid or expired join token",
		})
		return
	}

	if !roomExists(rname) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "requested room does not exist",
		})
		return
	}

	key := rname + ":" + name

	disconnectMu.Lock()
	timer, isReconnecting := disconnectTimers[key]
	if isReconnecting && timer != nil {
		timer.Stop()
		delete(disconnectTimers, key)
	}
	disconnectMu.Unlock()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &model.Client{
		Conn:     conn,
		Username: name,
	}

	addClientToRoom(rname, client)

	// Stream past chat history to the newly connected user
	if history, err := db.GetRoomHistory(rname, 50); err == nil {
		for _, hMsg := range history {
			_ = client.Conn.WriteJSON(hMsg)
		}
	}

	if !isReconnecting {
		msg := model.Message{
			Type:     1,
			Message:  "joined",
			Username: name,
			Room:     rname,
		}
		register <- msg
	}

	client.ReadMessageByRoom(name, rname, broadcast, unregister, func() {
		disconnectMu.Lock()
		if existing, ok := disconnectTimers[key]; ok && existing != nil {
			existing.Stop()
		}

		disconnectTimers[key] = time.AfterFunc(disconnectGracePeriod, func() {
			disconnectMu.Lock()
			delete(disconnectTimers, key)
			disconnectMu.Unlock()

			removeClientFromRoom(rname, client)
			msg := model.Message{
				Type:     2,
				Username: name,
				Message:  "left",
				Room:     rname,
			}
			unregister <- msg
		})
		disconnectMu.Unlock()
	})
}

func createRoom(name string, password string) bool {
	roomMu.Lock()
	defer roomMu.Unlock()

	if _, ok := rooms[name]; ok {
		return false
	}

	// Check MongoDB
	if existing, _ := db.GetRoomByName(name); existing != nil {
		rooms[name] = existing.Rpass
		return false
	}

	newRoom := model.Room{
		Rname:     name,
		Rpass:     password,
		CreatedAt: time.Now(),
	}

	if err := db.CreateRoom(newRoom); err != nil {
		return false
	}

	rooms[name] = password
	return true
}

func validateRoomPassword(name string, password string) (bool, bool) {
	roomMu.RLock()
	roomPassword, ok := rooms[name]
	roomMu.RUnlock()

	if ok {
		return true, password == roomPassword
	}

	// Fallback query from MongoDB
	room, err := db.GetRoomByName(name)
	if err != nil || room == nil {
		return false, false
	}

	roomMu.Lock()
	rooms[name] = room.Rpass
	roomMu.Unlock()

	return true, password == room.Rpass
}

func roomExists(name string) bool {
	roomMu.RLock()
	_, ok := rooms[name]
	roomMu.RUnlock()

	if ok {
		return true
	}

	room, err := db.GetRoomByName(name)
	if err != nil || room == nil {
		return false
	}

	roomMu.Lock()
	rooms[name] = room.Rpass
	roomMu.Unlock()

	return true
}

func addClientToRoom(room string, client *model.Client) {
	roomMu.Lock()
	defer roomMu.Unlock()

	if clientsroom[room] == nil {
		clientsroom[room] = make(map[*model.Client]bool)
	}

	clientsroom[room][client] = true
}

func removeClientFromRoom(room string, client *model.Client) {
	roomMu.Lock()
	defer roomMu.Unlock()

	delete(clientsroom[room], client)
	if len(clientsroom[room]) == 0 {
		delete(clientsroom, room)
		delete(rooms, room)
	}
}

func clientsInRoom(room string) []*model.Client {
	roomMu.RLock()
	defer roomMu.RUnlock()

	clients := make([]*model.Client, 0, len(clientsroom[room]))
	for client := range clientsroom[room] {
		clients = append(clients, client)
	}

	return clients
}

func newJoinToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func consumeJoinToken(token string, room string) bool {
	if token == "" || room == "" {
		return false
	}

	joinTokensMu.Lock()
	defer joinTokensMu.Unlock()

	claim, ok := joinTokens[token]
	if !ok {
		return false
	}

	delete(joinTokens, token)

	return claim.Room == room && time.Now().Before(claim.ExpiresAt)
}

func HandleMessagesByRoom() {
	for {
		select {
		case msg := <-register:
			go db.SaveMessage(msg)
			for _, client := range clientsInRoom(msg.Room) {
				err := client.Conn.WriteJSON(msg)
				if err != nil {
					continue
				}
			}

		case msg := <-unregister:
			go db.SaveMessage(msg)
			for _, client := range clientsInRoom(msg.Room) {
				err := client.Conn.WriteJSON(msg)
				if err != nil {
					continue
				}
			}

		case msg := <-broadcast:
			go db.SaveMessage(msg)
			for _, client := range clientsInRoom(msg.Room) {
				err := client.Conn.WriteJSON(msg)
				if err != nil {
					continue
				}
			}
		}
	}
}
