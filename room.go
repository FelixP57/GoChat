package main

// A room represents a discussion between two client
type Room struct {
	hub *Hub

	capacity int

	// Room clients
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
		capacity: 10,
		broadcast:	make(chan Event),
		clients:	make(map[*Client]bool),
		register:	make(chan *Client),
		unregister:	make(chan *Client),
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				if len(r.clients) == 0 {
					close(r.broadcast)
					close(r.register)
					close(r.unregister)
					return
				}
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
