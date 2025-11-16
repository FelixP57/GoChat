var conn;
var selectedRoom;
var username;

/*
* Event wraps all messages sent and received on the websocket
* the type is used as a RPC
*/
class Event {
    // type required but not payload
    constructor(type, payload) {
	this.type = type;
	this.payload = payload;
    }
}

/*
* SendMessageEvent sends messages to other clients
*/
class SendMessageEvent {
    constructor(message, room_id) {
	this.message = message;
	this.room_id = room_id;
    }
}

/*
* NewMessageEvent is messages from clients
*/
class NewMessageEvent {
    constructor(message, from, sent, room_id) {
	this.message = message;
	this.from = from;
	this.sent = sent;
	this.room_id = room_id;
    }
}

/*
* NewRoomEvent is for user rooms
*/
class NewRoomEvent {
    constructor(id, name, last_message) {
	this.id = id;
	this.name = name;
	this.last_message = last_message;
    }
}

/*
* CreateRoomEvent is for sending room creation request
*/
class CreateRoomEvent {
    constructor(username) {
	this.username = username
    }
}

/*
* routeEvent routes events to their handler based on type
*/
function routeEvent(event) {
    if (event.type == undefined) {
	alert("no 'type' field in event");
    }
    switch (event.type) {
	case "new_message":
	    // format payload
	    const messageEvent = Object.assign(new NewMessageEvent, event.payload);
	    setLastRoomMessage(messageEvent.room_id, messageEvent);
	    if (messageEvent.room_id != selectedRoom) {
		changeChatRoom(messageEvent.room_id);
	    } else {
		appendChatMessage(messageEvent);
	    }
	    break;
	case "new_room":
	    const roomEvent = Object.assign(new NewRoomEvent, event.payload);
	    appendRoomButton(roomEvent);
	    if (roomEvent.last_message.room_id == roomEvent.id) {
		setLastRoomMessage(roomEvent.id, roomEvent.last_message);
	    }
	    break;
	default:
	    alert("unsupported message type");
	    break;
	}
}

/*
* appendChatMessage takes in new messages and adds them to the chat
*/
function appendChatMessage(messageEvent) {
    var date = new Date(messageEvent.sent);
    const formattedMsg = `${date.toLocaleString()} ${messageEvent.from}: ${messageEvent.message}`;
    textarea = document.getElementById("chatmessages");
    textarea.innerHTML = textarea.innerHTML + "\n" + formattedMsg;
    textarea.scollTop = textarea.scrollHeight;
}

/*
* appendRoomButton takes in a new room and adds a button to join it
*/
function appendRoomButton(roomEvent) {
    roomDivs = document.getElementById("chatroom-selection");

    roomDiv = `<div class="room" id="room${roomEvent.id}" data-name="${roomEvent.name}">
	    <p class="room-name">${roomEvent.name}</p>
	    <p class="last-message"></p>
	</div>`

    roomDivs.innerHTML = roomDivs.innerHTML + roomDiv
}

/*
* changeChatRoom will update the value of selectedchat
* and also notify the server that it changes chatroom
*/
function changeChatRoom(roomId, roomName) {
    if (selectedRoom != 0) {
	document.getElementById("room" + selectedRoom.toString()).classList.remove("selected");
    }
    selectedRoom = roomId;
    document.getElementById("room" + selectedRoom.toString()).classList.add("selected");
    header = document.getElementById("chat-header");
    header.innerHTML = "Currently in chat: " + roomName;
    textarea = document.getElementById("chatmessages")
    textarea.innerHTML = "";
    sendEvent("get_messages", {"room_id": roomId});

    return false;
}

/*
* setLastRoomMessage sets the last message of a room in the
* appropriate html element
*/
function setLastRoomMessage(roomId, message) {
    p = document.querySelector(`#room${roomId} > .last-message`);
    p.innerHTML = `${message.from}: ${message.message}`
}
/*
* sendMessage sends a new message on the websocket
*/
function sendMessage() {
    var newMessage = document.getElementById("message");
    if (newMessage != null && selectedRoom != 0) {
	let outgoingEvent = new SendMessageEvent(newMessage.value, selectedRoom);
	sendEvent("send_message", outgoingEvent);
    } else {
	console.log("enter a message while in a room");
    }
    return false;
}

/*
* createRoom sends a room creation requests
*/
function createRoom() {
    var username = document.getElementById("with-username");
    if (username != null) {
	let outgoingEvent = new CreateRoomEvent(username.value);
	sendEvent("create_room", outgoingEvent);
    }
    return false;
}

/*
* sends the event on the websocket
*/
function sendEvent(eventName, payload) {
    const event = new Event(eventName, payload);
    if (conn != null) {
	conn.send(JSON.stringify(event));
    } else {
	alert("not authenticated");
    }
}

/*
* sends a login request to the server before connecting websocket
*/
function login() {
    let uname = document.getElementById("username").value;
    let formData = {
	"username": uname,
	"password": document.getElementById("password").value
    }
    // Send the request
    fetch("login", {
	method: 'post',
	body: JSON.stringify(formData),
	mode: 'cors',
    }).then((response) => {
	if (response.ok) {
	    return response.json();
	} else {
	    throw 'unauthorized';
	}
    }).then((data) => {
	// we have an OTP, request a connection
	connectWebsocket(data.token, uname);
    }).catch((e) => { alert(e) });
    return false
}

/* 
* sends a signup request to the server before connecting websocket
*/
function signup() {
    let uname = document.getElementById("username").value;
    let formData = {
	"username": uname,
	"password": document.getElementById("password").value
    }
    // Send the request
    fetch("signup", {
	method: 'post',
	body: JSON.stringify(formData),
	mode: 'cors',
    }).then((response) => {
	if (response.ok) {
	    return response.json();
	} else {
	    throw 'username already taken';
	}
    }).then((data) => {
	username = formData["username"]
	// we have an OTP, request a connection
	connectWebsocket(data.token, uname);
    }).catch((e) => { alert(e) });
    return false
}

/*
* logs out of the app
*/
function disconnect() {
    sendEvent("disconnect", {})
    document.getElementById("chatroom-selection").innerHTML = "";
    return false;
}

/*
* connects to websocket and adds listeners
*/
function connectWebsocket(token, uname) {
    if (window["WebSocket"]) {
	if (conn != null) {
	    disconnect();
	}
	conn = new WebSocket("wss://" + document.location.host + "/ws?token=" + token);

	// Onopen
	conn.onopen = function (evt) {
	    document.getElementById("chatmessages").innerHTML = "";
	    document.getElementById("chatroom-selection").innerHTML = "";
	    document.getElementById("messages").style.display = 'block';
	    document.getElementById("sidebar").style.display = 'block';
	    document.getElementById("profile-username").innerHTML = uname
	    document.getElementById("main").style.marginLeft = '30%';
	    document.getElementById("auth").style.display = 'none';
	    sendEvent("get_rooms", {});
	}

	conn.onclose = function (evt) {
	    document.getElementById("messages").style.display = 'none';
	    document.getElementById("sidebar").style.display = 'none';
	    document.getElementById("profile-username").innerHTML = '';
	    document.getElementById("main").style.marginLeft = '0';
	    document.getElementById("auth").style.display = 'block';
	}

	conn.onmessage = function (evt) {
	    // parse message as JSON
	    const eventData = JSON.parse(evt.data);
	    // assign JSON data to new Event
	    const event = Object.assign(new Event, eventData);
	    routeEvent(event);
	}
	
    } else {
	alert("Not supporting websockets");
    }
}

// once the website loads
window.onload = function () {
    // apply the listeners to the submit events to avoid redirects
    document.getElementById("create-chatroom").onsubmit = createRoom;
    document.getElementById("chatroom-message").onsubmit = sendMessage;
    document.getElementById("login-form").onsubmit = login;
    document.getElementById("signup-form").onsubmit = signup;
    document.getElementById("logout").onsubmit = disconnect;

    document.addEventListener("click", (e) => {
	room = e.target.closest(".room");
	if (room) {
	    changeChatRoom(parseInt(room.id.replace("room", "")), room.getAttribute("data-name"));
	}
    });

    conn = null;
    selectedRoom = 0;
};
