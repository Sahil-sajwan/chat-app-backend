package main

import (
	handler "chatapp/Handler"
	"chatapp/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.Use(middleware.OptionsMiddleware())
	r.POST("/create-room/:name", handler.CreateRoomHandler)
	r.POST("join-room/auth", handler.JoinRoomAuthHandler)
	r.GET("join-room/:name", handler.JoinRoomHandler)
	go handler.HandleMessagesByRoom()

	r.Run(":8080")

}
