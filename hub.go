package main

import (
	"log"
	"net/http"
	"errors"
	"context"
	"encoding/json"

	"github.com/gorilla/websocket"
)

var (
	/*
	* upgrader upgrades HTTP requests to persistent websocket connection
	*/
	upgrader = websocket.Upgrader{
		CheckOrigin:	 checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

var (
	newline	= []byte{'\n'}
	space	= []byte{' '}
)

var users map[string]string = make(map[string]string)

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case "https://localhost:8080":
		return true
	default:
		return false
	}
}

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Non full rooms
	roomsWaiting []*Room

	// Registered clients and their associated room
	clients map[*Client]bool

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// handlers -> functions that handle Events
	handlers map[string]EventHandler
}

func newHub(ctx context.Context) *Hub {
	h := &Hub{
		clients: 	make(map[*Client]bool),
		register:	make(chan *Client),
		unregister:	make(chan *Client),
		handlers: 	make(map[string]EventHandler),
	}
	h.setupEventHandlers()

	return h
}

// configures and adds all handlers
func (h *Hub) setupEventHandlers() {
	h.handlers[EventSendMessage] = SendMessageHandler
	h.handlers[EventChangeRoom] = ChatRoomHandler
}

// makes sure the events are handlers are correctly associated
func (h *Hub) routeEvent(event Event, c *Client) error {
	if handler, ok := h.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

func (h *Hub) signupHandler(w http.ResponseWriter, r *http.Request) {
	type userSignupRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req userSignupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	if _, ok := users[req.Username]; ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	users[req.Username] = req.Password

	type response struct {
		Token string `json:"token"`
	}

	token, err := generateJWT(req.Username)
	if err != nil {
		log.Println(err)
		return
	}

	resp := response{
		Token: token,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// verifies user authentification and returns a one time password
func (h *Hub) loginHandler(w http.ResponseWriter, r *http.Request) {
	type userLoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req userLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// authenticate user
	if password, ok := users[req.Username]; ok && password == req.Password {
		type response struct {
			Token string `json:"token"`
		}

		token, err := generateJWT(req.Username)
		if err != nil {
			log.Println(err)
			return
		}

		resp := response{
			Token: token,
		}

		data, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}

	// auth failure
	w.WriteHeader(http.StatusUnauthorized)
}

// serveWs handles websocket requests from the peer
func (h *Hub) serveWs(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		// user not authorized
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	
	claims, err := verifyJWT(token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	username, err := claims.GetSubject()
	if err != nil {
		log.Println(err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(h, conn, username)
	h.register <- client

	// allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writeMessages()
	go client.readMessages()
}

// add client to the clients list
func (h *Hub) addClient(client *Client) {
	h.clients[client] = true
	client.online = true
}

// remove client from clients list and end connection
func (h *Hub) removeClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		client.room.unregister <- client
		client.online = false
		client.conn.Close()
		delete(h.clients, client)
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
			if len(h.roomsWaiting) > 0 {
				room := h.roomsWaiting[0]
				if room.capacity == len(room.clients) + 1 {
					h.roomsWaiting = h.roomsWaiting[1:]
				}
				client.room = room
				room.register <- client

			} else {
				room := newRoom(h)
				go room.run()
				if room.capacity > 1 {
					h.roomsWaiting = append(h.roomsWaiting, room)
				}
				client.room = room
				room.register <- client
			}
		case client := <-h.unregister:
			h.removeClient(client)
		}
	}
			
}

