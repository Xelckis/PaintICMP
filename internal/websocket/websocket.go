package websocket

import (
	"fmt"
	"net/http"
	"sync"

	"paint/internal/common"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients   map[*Client]bool
	Mu        sync.Mutex
	Broadcast chan *[]common.Pixel
}

type Client struct {
	clientConn *websocket.Conn
	send       chan []byte
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
		Broadcast: make(chan *[]common.Pixel, 100),
	}
}

func (h *Hub) WsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	client := &Client{clientConn: conn, send: make(chan []byte, 100)}
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
		pixelSlice := <-h.Broadcast

		pixelBytes := make([]byte, 0, len(*pixelSlice)*3)
		for _, p := range *pixelSlice {
			pixelBytes = append(pixelBytes, p.X, p.Y, p.Color)
		}

		common.PoolPut(pixelSlice)

		h.Mu.Lock()
		for client := range h.Clients {
			select {
			case client.send <- pixelBytes:
			default:
				client.clientConn.Close()
			}
		}
		h.Mu.Unlock()
	}
}

func writePump(conn *Client) {
	defer conn.clientConn.Close()

	for pixelBytes := range conn.send {
		if err := conn.clientConn.WriteMessage(websocket.BinaryMessage, pixelBytes); err != nil {
			fmt.Println("Error writing binary:", err)
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
