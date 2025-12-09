package main

import (
	"paint/internal/icmp"
	"paint/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	go icmp.FilterICMP()
	router := gin.Default()
	router.GET("/ws", websocket.WsHandler)
	router.Run()
}
