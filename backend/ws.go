package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	db *Db

	conn *websocket.Conn

	send chan []byte
}

type Db struct {
	clients map[*Client]bool

	messages [][]byte

	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func DbNew() *Db {
	return &Db{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (db *Db) run() {
	for {
		select {
		case client := <-db.register:
			db.clients[client] = true
			if len(db.messages) > 0 {
				messages, _ := json.Marshal(db.messages)
				client.send <- messages
			}
		case client := <-db.unregister:
			log.Println("Unregistered new client ", client)
			if _, ok := db.clients[client]; ok {
				delete(db.clients, client)
				close(client.send)
			}
		case message := <-db.broadcast:
			db.messages = append(db.messages, message)
			response, _ := json.Marshal(db.messages)
			for client := range db.clients {
				select {
				case client.send <- response:
				default:
					close(client.send)
				}
			}
		}
	}
}

func (c *Client) readListener() {
	defer func() {
		c.db.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(message)
		c.db.broadcast <- message
	}
}

func (c *Client) writeListener() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func ServeWs(db *Db, w http.ResponseWriter, r *http.Request) {
	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}).Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{db: db, conn: conn, send: make(chan []byte, 256)}
	client.db.register <- client

	go client.writeListener()
	go client.readListener()
}
