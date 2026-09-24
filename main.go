package main

import (
	handler "chatapp/Handler"
	"chatapp/db"
	"chatapp/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.InitMongo(); err != nil {
		log.Printf("MongoDB initialization notice: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.OptionsMiddleware())
	r.POST("/create-room/:name", handler.CreateRoomHandler)
	r.POST("/join-room/auth", handler.JoinRoomAuthHandler)
	r.GET("/join-room/:name", handler.JoinRoomHandler)
	go handler.HandleMessagesByRoom()

	r.Run(":8080")
}
