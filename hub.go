package main

import (
	"log"
	"net/http"
	"errors"
	"context"
	"encoding/json"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
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
	roomsWaiting []int

	// all rooms
	rooms map[int]*Room

	// Registered clients and their associated user
	clients map[*Client]*User

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// handlers -> functions that handle Events
	handlers map[string]EventHandler

	db *Database
}

func newHub(ctx context.Context) *Hub {
	h := &Hub{
		clients: 	make(map[*Client]*User),
		rooms:		make(map[int]*Room),
		register:	make(chan *Client),
		unregister:	make(chan *Client),
		handlers: 	make(map[string]EventHandler),
	}
	h.db = getDb(h)
	h.setupEventHandlers()
	h.loadRooms()

	return h
}

func (h *Hub) loadRooms() {
	rooms, err := h.db.getRooms()
	if err != nil {
		log.Println(err)
	}
	h.rooms = rooms
	for i := range rooms {
		go rooms[i].run()
	}
}

// configures and adds all handlers
func (h *Hub) setupEventHandlers() {
	h.handlers[EventSendMessage] = SendMessageHandler
	h.handlers[EventChangeRoom] = ChatRoomHandler
	h.handlers[EventDisconnectClient] = DisconnectClientHandler
	h.handlers[EventGetMessages] = GetMessagesHandler
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


	if _, err := h.db.getUserByUsername(req.Username); err == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost) 

	h.db.addUser(newUser(req.Username, string(bytes))) 

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
	if user, err := h.db.getUserByUsername(req.Username); err == nil && bcrypt.CompareHashAndPassword([]byte(user.password), []byte(req.Password)) == nil {
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

	user, err := h.db.getUserByUsername(username)
	if err != nil  {
		log.Println(err)
	}
	client := newClient(h, conn, user)
	h.register <- client

	// allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writeMessages()
	go client.readMessages()
}

// add client to the clients list
func (h *Hub) addClient(client *Client) {
	h.clients[client] = client.user
}

// remove client from clients list and end connection
func (h *Hub) removeClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		room, ok := h.rooms[client.user.roomId]
		if ok {
			room.unregister <- client
		}
		client.conn.Close()
		delete(h.clients, client)
	}
}

func (h *Hub) run() {
	defer h.db.closeDb()
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
			if client.user.roomId != 0 {
				room := h.rooms[client.user.roomId]
				room.register <- client
				continue
			}
			if len(h.roomsWaiting) > 0 {
				roomId := h.roomsWaiting[0]
				room, ok := h.rooms[roomId]
				if !ok {
					log.Println("Unknown room")
				}
				if room.capacity == len(room.users) + 1 {
					h.roomsWaiting = h.roomsWaiting[1:]
				}
				client.user.roomId = roomId
				h.db.addUserToRoom(client.user.username, roomId)
				h.db.updateUserRoom(client.user, roomId)
				room.register <- client

			} else {
				room := newRoom(h)
				roomId := h.db.addRoom(room)
				h.rooms[roomId] = room
				go room.run()
				if room.capacity > 1 {
					h.roomsWaiting = append(h.roomsWaiting, roomId)
				}
				client.user.roomId = roomId
				h.db.addUserToRoom(client.user.username, roomId)
				h.db.updateUserRoom(client.user, roomId)
				room.register <- client
			}
		case client := <-h.unregister:
			h.removeClient(client)
		}
	}
			
}

