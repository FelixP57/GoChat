package main

import (
	"log"
	"time"
	"encoding/json"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 
)


// Client is a the intermediary between the websocket connection and the hub
type Client struct {
	hub *Hub

	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound message
	send chan Event

	user *User

	online bool
}

func newClient(h *Hub, conn *websocket.Conn, user *User) *Client {
	return &Client{
		hub: h,
		conn: conn,
		user: user,
		send: make(chan Event),
	}
}

// reads messages from the websocket connection to the hub
// ran in a goroutine for each connection, so that there can only be one read at a time
func (c *Client) readMessages() {
	defer func() {
		c.hub.unregister <- c
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil})
	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// unmarshal incoming data into Event
		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error unmarshalling message: %v", err)
			break
		}

		// route Event
		if err := c.hub.routeEvent(request, c); err != nil {
			log.Println("error handling message: ", err)
		}
	}
}

// writes messages from the hub to the websocket connection.
// one writeMessages goroutine is started with each connection to ensure only one write at a time
func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.hub.unregister <- c
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err)
				return
			}
			w.Write(data)

			// Add queued chat messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				data, err := json.Marshal(message)
				if err != nil {
					log.Println(err)
					return
				}

				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				w.Write(data)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

