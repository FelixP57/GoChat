import { Event, 
    SendMessageEvent, 
    NewMessageEvent, 
    NewRoomEvent, 
    CreateRoomEvent  
} from './events.js';

// const API_DOMAIN = import.meta.env.VITE_API_DOMAIN
const API_DOMAIN="localhost:8080";

var conn = null;

function sendEvent(eventName, payload) {
    const event = new Event(eventName, payload);
    console.log(event);
    if (conn != null) {
	conn.send(JSON.stringify(event));
    } else {
	alert("not authenticated");
    }
}

function connectWebsocket(token, username, callback) {
    if (window["WebSocket"]) {
	if (conn != null) {
	    sendEvent("disconnect", {})
	}

	console.log("Attempting connection...");
	conn = new WebSocket("wss://" + API_DOMAIN + "/ws?token=" + token);

	// Onopen
	conn.onopen = function (evt) {
	    console.log("Successfully connected");
	    sendEvent("get_rooms", {});
	}

	conn.onclose = function (evt) {
	    console.log("Socket closed connection", event);
	}

	conn.onmessage = function (evt) {
	    // parse message as JSON
	    const eventData = JSON.parse(evt.data);
	    // assign JSON data to new Event
	    const event = Object.assign(new Event, eventData);
	    callback(event);
	}
	    
	conn.onerror = error => {
	    console.log("Socket error: ", error);
	};
	
    } else {
	alert("Not supporting websockets");
    }
}

export { 
    connectWebsocket, 
    sendEvent, 
    API_DOMAIN,
    SendMessageEvent, 
    NewMessageEvent, 
    NewRoomEvent, 
    CreateRoomEvent,
};
