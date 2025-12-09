package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Pixel struct {
	X, Y, Color string
}

var (
	PixelChan = make(chan Pixel)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	handleConnection(conn)
}

func handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Printf("Cliente desconectado (read error): %v", err)
				break
			}
		}
	}()
	for {
		pixel := <-PixelChan
		log.Printf("Recebi o pixel aqui do outro lado: X: %s Y: %s Color: %s", pixel.X, pixel.Y, pixel.Color)
		if err := conn.WriteJSON(pixel); err != nil {
			fmt.Println("Error writing JSON:", err)
			break
		}
	}
}
