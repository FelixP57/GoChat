package main

// A room represents a discussion between two client
type Room struct {
	hub *Hub

	id int

	capacity int

	name string

	// Authorized users
	users map[*User]bool

	// Connected clients
	clients map[*Client]bool
	
	// Inbound messages from the clients
	broadcast chan Event

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client
}

func newRoom(hub *Hub) *Room {
	return &Room{
		hub: hub,
		capacity: 2,
		name: "",
		broadcast:	make(chan Event),
		users:		make(map[*User]bool),
		clients:	make(map[*Client]bool),
		register:	make(chan *Client),
		unregister:	make(chan *Client),
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.users[client.user] = true
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
			}
		case event := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- event:
				default:
					r.hub.unregister <- client 	
				}
			}
		}
	}
}
