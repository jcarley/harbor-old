package service

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Hub struct {
	connections map[*Connection]bool
	register    chan *Connection
	unregister  chan *Connection
	refresh     chan bool
}

func NewHub() *Hub {
	return &Hub{
		refresh:     make(chan bool),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: make(map[*Connection]bool),
	}
}

func (h *Hub) Register(c *Connection) {
	h.register <- c
}

func (h *Hub) Unregister(c *Connection) {
	h.unregister <- c
}

func (h *Hub) Refresh() {
	h.refresh <- true
}

func (h *Hub) onRefresh() []byte {
	// sort.Sort(Reverse{h.runList})
	// recent := h.runList.GetRecent(0, 10)
	// bytes, err := json.Marshal(recent)
	var err error = nil
	bytes := []byte{}
	if err != nil {
		panic(err.Error())
	}
	return bytes
}

func (h *Hub) HubLoop() {
	for {
		select {
		case c := <-h.register:
			fmt.Println("Connect")
			h.connections[c] = true
			bytes := h.onRefresh()
			c.send <- bytes
		case c := <-h.unregister:
			fmt.Println("Disconnect")
			delete(h.connections, c)
			close(c.send)
		case <-h.refresh:
			fmt.Println("Refreshing")
			bytes := h.onRefresh()
			for c := range h.connections {
				select {
				case c.send <- bytes:
				default:
					delete(h.connections, c)
					close(c.send)
					go c.ws.Close()
				}
			}
		}
	}
}

type Connection struct {
	ws   *websocket.Conn
	send chan []byte
}

func NewConnection(ws *websocket.Conn) *Connection {
	return &Connection{send: make(chan []byte, 256), ws: ws}
}

func (c *Connection) Reader() {
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			fmt.Printf("Error in websocket read: %s\n", err.Error())
			break
		}
		fmt.Printf("Message received: %s\n", msg)
		// TODO: Do something with the msg
	}
	c.ws.Close()
}

func (c *Connection) Writer() {
	for msg := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Printf("Error in websocket write: %s\n", err.Error())
			break
		}
	}
	c.ws.Close()
}
