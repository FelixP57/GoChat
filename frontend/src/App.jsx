import React, { Component } from "react";
import './App.css';
import { 
    connectWebsocket, 
    sendEvent, 
    API_DOMAIN,     
    SendMessageEvent, 
    NewMessageEvent, 
    NewRoomEvent, 
    CreateRoomEvent,
} from "./api";

import Header from './components/Header/Header';
import AuthForm from './components/AuthForm/AuthForm';
import Messages from './components/Messages/Messages';
import RoomList from './components/RoomList/RoomList';
import Profile from './components/Profile/Profile';
import RoomCreationForm from './components/RoomCreationForm/RoomCreationForm';

class App extends Component {
    constructor(props) {
	super(props);
	this.state = {
	    username: null,
	    selectedRoom: null,
	    messages: [],
	    rooms: new Map(),
	};
	this.login = this.login.bind(this);
	this.signup = this.signup.bind(this);
	this.routeEvent = this.routeEvent.bind(this);
	this.addRoom = this.addRoom.bind(this);
	this.changeChatRoom = this.changeChatRoom.bind(this);
	this.sendMessage = this.sendMessage.bind(this);
	this.setLastRoomMessage = this.setLastRoomMessage.bind(this);
	this.disconnect = this.disconnect.bind(this);
    }

    routeEvent(event) {
	if (event.type == undefined) {
	    alert("no 'type' field in event");
	}
	switch (event.type) {
	    case "new_message":
		// format payload
		const messageEvent = Object.assign(new NewMessageEvent, event.payload);
		this.setLastRoomMessage(messageEvent.room_id, messageEvent);
		if (!this.state.selectedRoom || messageEvent.room_id != this.state.selectedRoom.id) {
		    this.changeChatRoom(messageEvent.room_id);
		} else {
		    this.setState(prevState => ({
			messages: [...prevState.messages, messageEvent]
		    }));
		}
		break;
	    case "new_room":
		const roomEvent = Object.assign(new NewRoomEvent, event.payload);
		this.setState(prevState => ({
		    rooms: this.addRoom(roomEvent, prevState.rooms)
		}));
		break;
	    default:
		alert("unsupported message type");
		break;
	    }
    }

    setLastRoomMessage(roomId, message) {
	let rooms = new Map(this.state.rooms);
	rooms.get(parseInt(roomId)).last_message = message;
	this.setState({rooms: rooms});
    }

    addRoom(roomEvent, rooms) {
	rooms.set(parseInt(roomEvent.id), roomEvent);
	return rooms;
    }

    changeChatRoom(roomId) {
	this.setState({
	    selectedRoom: this.state.rooms.get(roomId),
	    messages: [],
	});
	sendEvent("get_messages", {"room_id": roomId});
	return false;
    }

    sendMessage(message) {
	let outgoingEvent = new SendMessageEvent(message, parseInt(this.state.selectedRoom.id));
	sendEvent("send_message", outgoingEvent);
	return false;
    }

    createRoom(username) {
	let outgoingEvent = new CreateRoomEvent(username);
	sendEvent("create_room", outgoingEvent);
	return false;
    }

    login(event) {
	event.preventDefault();
	let username = document.getElementById("username").value;
	let formData = {
	    "username": username,
	    "password": document.getElementById("password").value
	}
	// Send the request
	fetch(`https://${API_DOMAIN}/login`, {
	    method: 'post',
	    body: JSON.stringify(formData),
	    mode: 'cors',
	    headers: {
		"Content-Type": "application/json",
	    },
	}).then((response) => {
	    if (response.ok) {
		return response.json();
	    } else {
		throw 'unauthorized';
	    }
	}).then((data) => {
	    this.setState({username: username});
	    // we have an OTP, request a connection
	    connectWebsocket(data.token, username, this.routeEvent);
	}).catch((e) => { alert(e) });

	return false;
    }

    signup(event) {
	event.preventDefault();
	let username = document.getElementById("username").value;
	let formData = {
	    "username": username,
	    "password": document.getElementById("password").value
	}
	// Send the request
	fetch(`https://${API_DOMAIN}/signup`, {
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
	    this.setState({username: username});
	    // we have an OTP, request a connection
	    connectWebsocket(data.token, username, this.routeEvent);
	}).catch((e) => { alert(e) });

	return false;
    }

    disconnect() {
	sendEvent("disconnect", {})
	return false;
    }

    render() {
	return (
	    <div className="App">
		<Header />

		{this.state.username  && 
		<div id="sidebar">
		    <RoomList rooms={this.state.rooms} selectedRoom={this.state.selectedRoom} changeRoom={this.changeChatRoom} />
		   <br /> 
		    <RoomCreationForm createRoom={this.createRoom} />

		    <Profile username={this.state.username} logout={this.disconnect} />
		</div>
		}
	    
		<div 
		    className="main"
		    style={{marginLeft: this.state.username ? '200px' : 'auto'}}
		>

		    {this.state.username && this.state.selectedRoom &&
		    <Messages 
			messages={this.state.messages} 
			roomName={this.state.selectedRoom.name} 
			sendMessage={this.sendMessage}
		    />
		    }

		    {this.state.username == null && 
		    <AuthForm 
			loginHandler={this.login} 
			signupHandler={this.signup}
		    />
		    }

		</div>
	    </div>
	);
    }
}

export default App;
