package main

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	WRITE_WAIT = 10 * time.Second

	// Time allowed to readPump the next pong message from the peer.
	PONG_WAIT = 60 * time.Second

	// Maximum message size allowed from peer.
	MAX_MSG_SIZE = 512

	// Send pings to peer with this period. Must be less than PONG_WAIT.
	PING_PERIOD = (PONG_WAIT * 9) / 10
)

type FromMsg struct {
	From string
	Data interface{}
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
	id   string
}

func (client *Client) readPump() {
	defer func() {
		client.conn.Close()
		client.hub.unregister <- client.id
	}()

	client.conn.SetReadLimit(MAX_MSG_SIZE)
	client.conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
		return nil
	})
	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		client.hub.receive <- string(msg)
	}
}

func (client *Client) writePump() {
	ticker := time.NewTicker(PING_PERIOD)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
			if !ok { // the hub has closed the chan
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
			}

			err := client.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type Hub struct {
	connections map[string]*Client
	unregister  chan string
	receive     chan string
}

func newHub() *Hub {
	return &Hub{
		connections: make(map[string]*Client),
		unregister:  make(chan string),
		receive:     make(chan string),
	}
}

func (h *Hub) add(id string, c *websocket.Conn) error {
	if _, ok := h.connections[id]; ok {
		return errors.New("already connected w this id")
	}
	cn := &Client{
		conn: c,
		send: make(chan []byte, 3),
		hub:  h,
		id:   id,
	}
	h.connections[id] = cn
	go cn.readPump()
	go cn.writePump()
	return nil
}

func (h *Hub) sendAll(val []byte) {
	for id, client := range h.connections {
		select {
		case client.send <- val:
		default:
			h.end(id)
		}
	}
}

func (h *Hub) send(to string, val []byte) {
	if client, ok := h.connections[to]; ok {
		select {
		case client.send <- val:
		default:
			h.end(to)
		}
	}
}

func (h *Hub) end(id string) {
	if client, ok := h.connections[id]; ok {
		close(client.send)
		delete(h.connections, id)
	}
}
