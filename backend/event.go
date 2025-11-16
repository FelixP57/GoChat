package main

import (
	"fmt"
	"errors"
	"time"
	"encoding/json"
)

// Event is the messages sent over the websocket to distinguish different actions
type Event struct {
	Type 	string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// func signatureto affect messages on the socket based on type
type EventHandler func(event Event, c *Client) error

const (
	// event name for new chat messages sent
	EventSendMessage = "send_message"
	// response to send_message
	EventNewMessage = "new_message"
	// disconenct client
	EventDisconnectClient = "disconnect"
	// get room messages
	EventGetMessages = "get_messages"
	// get rooms
	EventGetRooms = "get_rooms"
	// new room
	EventNewRoom = "new_room"
	// create room
	EventCreateRoom = "create_room"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From 	string `json:"from"`
	RoomId  int    `json:"room_id"`
}

// returned when responding to send_message or get_messages
type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

// returned when responding to get_rooms
type NewRoomEvent struct {
	Id int `json:"id"`
	Name string `json:"name"`
	LastMessage NewMessageEvent `json:"last_message"`
}

type CreateRoomEvent struct {
	Username string `json:"username"`
}

type GetMessagesEvent struct {
	RoomId int `json:"room_id"`
}

func SendMessageHandler(event Event, c *Client) error {
	var chatevent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	var broadMessage NewMessageEvent

	fmt.Println(chatevent.Message)
	fmt.Println(chatevent.RoomId)
	broadMessage.Sent = time.Now()
	broadMessage.Message = chatevent.Message
	broadMessage.From = c.user.username
	broadMessage.RoomId = chatevent.RoomId

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	c.hub.db.addMessage(broadMessage, chatevent.RoomId)

	// place payload in an event
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewMessage

	room, ok := c.hub.rooms[broadMessage.RoomId]
	if !ok {
		fmt.Errorf("error retrieving room by id: %v", err)
	}

	room.lastMessage = broadMessage
	room.broadcast <- outgoingEvent

	return nil
}

func DisconnectClientHandler(event Event, c *Client) error {
	c.hub.removeClient(c)
	return nil
}

func GetMessagesHandler(event Event, c *Client) error {
	var e GetMessagesEvent;
	if err := json.Unmarshal(event.Payload, &e); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}
	events := c.hub.db.getMessages(e.RoomId)
	for i := range events {
		data, err := json.Marshal(events[i])
		if err != nil {
			return fmt.Errorf("failed to marshal broadcast message: %v", err)
		}

		// place payload in an event
		var outgoingEvent Event
		outgoingEvent.Payload = data
		outgoingEvent.Type = EventNewMessage

		c.send <- outgoingEvent
	}
	return nil
}

func GetRoomsHandler(event Event, c *Client) error {
	roomIds := c.hub.db.getRooms(c.user.username)
	for i := range roomIds {
		room := c.hub.rooms[roomIds[i]]
		roomName := room.name
		if roomName == "" {
			users := c.hub.rooms[roomIds[i]].users
			for username := range users {
				if (username != c.user.username) {
					roomName += username
					roomName += ", ";
				}
			}
			roomName = roomName[:len(roomName) - 2];
		}
		lastMessage := c.hub.db.getLastRoomMessage(roomIds[i])
		event := NewRoomEvent{Id: roomIds[i], Name: roomName, LastMessage: lastMessage}
		data, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal broadcast message: %v", err)
		}

		var outgoingEvent Event
		outgoingEvent.Payload = data
		outgoingEvent.Type = EventNewRoom

		c.send <- outgoingEvent
	}
	return nil
}

func CreateRoomHandler(event Event, c *Client) error {
	var createRoom CreateRoomEvent
	if err := json.Unmarshal(event.Payload, &createRoom); err != nil {
		fmt.Errorf("Error unmarshalling event: %v", err)
	}
	// check if other user exists
	user, err := c.hub.db.getUserByUsername(createRoom.Username)
	if err != nil {
		fmt.Errorf("user not found: %v", err)
	}
	
	// check if room already exists
	var roomEvent NewRoomEvent
	var room *Room
	id, err := c.hub.db.getRoomByUsers(c.user.username, user.username)
	if err != nil {
		if !errors.Is(err, RoomNotFoundError) {
			return err
		}
		// create room
		room = newRoom(c.hub)
		go room.run()
		id := c.hub.db.addRoom(room)
		
		room.register <- c.user
		room.register <- user
		c.hub.db.addUserToRoom(c.user.username, id)
		c.hub.db.addUserToRoom(user.username, id)
		c.hub.rooms[id] = room
		roomEvent = NewRoomEvent{Id: id, Name: user.username, LastMessage: NewMessageEvent{}}
	} else {
		room = c.hub.rooms[id]
		roomEvent = NewRoomEvent{Id: id, Name: user.username, LastMessage: room.lastMessage}
	}	

	// broadcast NewRoomEvent
	data, err := json.Marshal(roomEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewRoom
	room.broadcast <- outgoingEvent

	return nil
}

