package main

// A room represents a discussion between two client
type Room struct {
	hub *Hub

	id int

	capacity int

	name string

	// Authorized users' username
	users map[string]bool

	// Inbound messages from the clients
	broadcast chan Event

	// Register requests from the clients
	register chan *User

	// Unregister requests from clients
	unregister chan *User
}

func newRoom(hub *Hub) *Room {
	return &Room{
		hub: hub,
		capacity: 2,
		name: "",
		broadcast:	make(chan Event),
		users:		make(map[string]bool),
		register:	make(chan *User),
		unregister:	make(chan *User),
	}
}

func (r *Room) run() {
	for {
		select {
		case user := <-r.register:
			r.users[user.username] = true
		case user := <-r.unregister:
			if _, ok := r.users[user.username]; ok {
				delete(r.users, user.username)
			}
		case event := <-r.broadcast:
			for user := range r.users {
				for client := range r.hub.clients[user] {
					select {
					case client.send <- event:
					default:
						r.hub.unregister <- client 	
					}
				}
			}
		}
	}
}
