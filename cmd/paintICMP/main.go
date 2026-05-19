package main

import (
	"paint/internal/icmp"
	"paint/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	go websocket.GlobalHub.Start()
	go icmp.FilterICMP()
	router := gin.Default()
	router.GET("/ws", websocket.GlobalHub.WsHandler)
	router.GET("/", func(c *gin.Context) {
		c.File("web/index.html")
	})
	router.Run()
}
