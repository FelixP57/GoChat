import React, { Component } from "react";
import './Messages.scss';

class Messages extends Component {
    constructor(props) {
	super(props);
	this.getOnlineUsers = this.getOnlineUsers.bind(this);
	this.sendMessage = this.sendMessage.bind(this);
	this.formatDate = this.formatDate.bind(this);
	this.scrollToBottom = this.scrollToBottom.bind(this);

	this.state = {
	    message: '',
	    onlineUsers: this.getOnlineUsers(),
	};
    }

    scrollToBottom() {
	this.messagesEnd.scrollIntoView({block: "end"});
    }

    getOnlineUsers() {
	var onlineUsers = new Array();
	this.props.room.users.forEach((user, index) => {
	    if (user.online) {
		onlineUsers.push(user.username);
	    }
	});
	return onlineUsers;
    }

    componentDidMount() {
	this.scrollToBottom();
    }

    componentDidUpdate() {
	this.scrollToBottom();
    }

    formatDate(date) {
	date = new Date(date);
	return date.toLocaleString('default', {
	    "hour": "2-digit",
	    "minute": "2-digit",
	    "hourCycle": "h24",
	});
    }

    getDay(date) {
	date = new Date(date);
	return date.toLocaleString('default', {
	    "month": "long",
	    "day": "numeric",
	    "year": "numeric",
	});
    }
    
    sendMessage(event) {
	event.preventDefault();
	if (this.state.message != '') {
	    this.props.sendMessage(this.state.message);
	} else {
	    console.log("enter a message while in a room");
	}
	this.setState({message: ""});
	return false;
    }

    render() {
	return (
	    <div id="messages">
		<h3 id="chat-header">Currently in chat: {this.props.room.name}</h3>
		<h3>{this.state.onlineUsers.join(",")}</h3>

		<div className="messagearea" id="chatmessages">
		    {this.props.messages.map((message, index, array) => {
			let day = this.getDay(message.sent);
			let show = (index == 0 || day != this.getDay(array[index-1].sent));
			return (
			<div key={index}>
			    <p className="date-separator">
				{show && this.getDay(message.sent)}
			    </p>
			    <p className="message">
				{this.formatDate(message.sent)} {message.from}: {message.message}
			    </p>
			</div>
		    );})}
		    <div ref={(el) => {this.messagesEnd = el;}}></div>
		</div>

		<br />
		<form id="chatroom-message" onSubmit={this.sendMessage} >
		    <label htmlFor="message">Message:</label>
		    <input type="text" id="message" name="message" onChange={(e) => this.setState({message: e.target.value})} value={this.state.message} autoFocus /><br /><br />
		    <input className="submit" type="submit" value="Send message" />
		</form>
	    </div>
	);
    }
}

export default Messages;

