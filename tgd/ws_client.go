package tgd

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = ws.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}

func serverWs(chathub *ChatHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("Unable to upgrade the connection:", err)
	}

	client := Client{chathub, uuid.NewString(), conn, make(chan []byte, 256)}
	chathub.register <- &client
	go client.readPump()
	go client.writePump()
}

type Client struct {
	hub  *ChatHub
	id   string
	conn *ws.Conn
	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		fmt.Printf("message from client %s: %s\n", c.id, message)
	}

	c.hub.unregister <- c
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}

		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}

			wsWriter, err := c.conn.NextWriter(ws.TextMessage)
			if err != nil {
				log.Printf("Error while getting next writer for %q.\n%v\n", c.id, err)
				return
			}
			wsWriter.Write(msg)

			n := len(c.send)
			for range n {
				wsWriter.Write(newline)
				wsWriter.Write(<-c.send)
			}

			if err = wsWriter.Close(); err != nil {
				log.Printf("Writing failed for %q.\n%v\n", c.id, err)
				return
			}
		}
	}
}
