import React, { Component } from "react";
import './Messages.scss';

class Messages extends Component {
    constructor(props) {
	super(props);
	this.state = {
	    message: '',
	};
	this.sendMessage = this.sendMessage.bind(this);
	this.formatDate = this.formatDate.bind(this);
    }

    formatDate(date) {
	date = new Date(date);
	const daysDiff = Math.floor((new Date() - date)/(1000*60*60*24));
	if (daysDiff == 0) {
	    return date.toLocaleString('default', {
		"hour": "2-digit",
		"minute": "2-digit",
		"hourCycle": "h24",
	    });
	} else {
	    return date.toLocaleString('default', {
		"hour": "2-digit",
		"minute": "2-digit",
		"hourCycle": "h24",
		"day": "numeric",
		"month": "long",
	    });
	}
    }
    
    sendMessage(event) {
	event.preventDefault();
	if (this.state.message != '') {
	    this.props.sendMessage(this.state.message);
	} else {
	    console.log("enter a message while in a room");
	}
	return false;
    }

    render() {
	return (
	    <div id="messages">
		<h3 id="chat-header">Currently in chat: {this.props.roomName}</h3>

		<div className="messagearea" id="chatmessages">
		    {this.props.messages.map((message, index) => (
			<p key={index}>
			    {this.formatDate(message.sent)} {message.from}: {message.message}
			</p>
		    ))}
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

