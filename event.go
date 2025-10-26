package main

import (
	"fmt"
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
	// switch rooms
	EventChangeRoom = "change_room"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From 	string `json:"from"`
}

// returned when responding to send_message
type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

func SendMessageHandler(event Event, c *Client) error {
	var chatevent SendMessageEvent;
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	var broadMessage NewMessageEvent

	broadMessage.Sent = time.Now()
	broadMessage.Message = chatevent.Message
	broadMessage.From = c.username

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// place payload in an event
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewMessage

	c.room.broadcast <- outgoingEvent
	
	return nil
}

type ChangeRoomEvent struct {
	Name string `json:"name"`
}

// handles switching of chatrooms
func ChatRoomHandler(event Event, c *Client) error {
	var changeRoomEvent ChangeRoomEvent;
	if err := json.Unmarshal(event.Payload, &changeRoomEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	// c.chatroom = changeRoomEvent.Name

	return nil
}

