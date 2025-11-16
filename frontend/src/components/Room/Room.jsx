import React, { Component } from "react";
import './Room.scss';

class Room extends Component {
    constructor(props) {
	super(props);
	this.roomId = this.props.roomId;
	this.roomName = this.props.roomName;
	this.lastMessage = this.props.lastMessage;
	this.changeRoom = this.changeRoom.bind(this);
	this.formatDate = this.formatDate.bind(this);
    }

    changeRoom() {
	this.props.changeRoom(this.roomId);
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
	} else if (daysDiff < 7) {
	    return `${daysDiff} d.`;
	} else {
	    return date.toLocaleString('default', {
		"day": "numeric",
		"month": "short",
	    });
	}
    }

    render() {
	return (
	    <div className="room" id={'room' + this.roomId} onClick={this.changeRoom}>
		<p className="room-name">{this.roomName}</p>
		<div className="last-message">
		    <p className="last-message-author">{this.lastMessage && this.lastMessage.from}: </p>
		    <p className="last-message-message">{this.lastMessage && this.lastMessage.message}</p>
		    <p className="last-message-date">{this.lastMessage && this.formatDate(this.lastMessage.sent)}</p>
		</div>
	    </div>
	);
    }
}

export default Room
