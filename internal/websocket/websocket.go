package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Pixel struct {
	X, Y, Color string
}

type Hub struct {
	Clients   map[*Client]bool
	Mu        sync.Mutex
	Broadcast chan Pixel
}

type Client struct {
	clientConn *websocket.Conn
	send       chan Pixel
}

var (
	GlobalHub = NewHub()
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func NewHub() *Hub {
	return &Hub{
		Clients:   make(map[*Client]bool),
		Broadcast: make(chan Pixel, 10000),
	}
}

func (h *Hub) WsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	client := &Client{clientConn: conn, send: make(chan Pixel, 10000)}
	h.AddClient(client)

	go writePump(client)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.RemoveClient(client)
			break
		}
	}

}

func (h *Hub) Start() {
	for {
		pixel := <-h.Broadcast

		h.Mu.Lock()
		for client := range h.Clients {
			select {
			case client.send <- pixel:
			default:
				client.clientConn.Close()
			}
		}
		h.Mu.Unlock()
	}
}

func writePump(conn *Client) {
	defer conn.clientConn.Close()

	for pixel := range conn.send {
		log.Printf("Recebi o pixel aqui do outro lado: X: %s Y: %s Color: %s", pixel.X, pixel.Y, pixel.Color)
		if err := conn.clientConn.WriteJSON(pixel); err != nil {
			fmt.Println("Error writing JSON:", err)
			return
		}
	}
}

func (h *Hub) AddClient(conn *Client) {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	h.Clients[conn] = true
	fmt.Println("Novo cliente conectado:", conn.clientConn.RemoteAddr().String())
}

func (h *Hub) RemoveClient(conn *Client) {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	delete(h.Clients, conn)
	close(conn.send)
	fmt.Println("Cliente desconectado:", conn.clientConn.RemoteAddr().String())
}
