import React, { Component } from "react";
import './Room.scss';

class Room extends Component {
    constructor(props) {
	super(props);
	this.changeRoom = this.changeRoom.bind(this);
	this.formatDate = this.formatDate.bind(this);
    }

    changeRoom() {
	this.props.changeRoom(this.props.room.id);
    }

    formatDate(date) {
	date = new Date(date);
	let today = new Date();
	today.setHours(0,0,0,0);
	const daysDiff = Math.ceil((today - date)/(1000*60*60*24));
	if (daysDiff <= 0) {
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
	    <div className="room" id={'room' + this.props.room.id} onClick={this.changeRoom}>
		<p className="room-name">{this.props.room.name}</p>
		<div className="last-message">
		    <p className="last-message-author">{this.props.room.last_message.from && this.props.room.last_message.from + ": "}</p>
		    <p className="last-message-message">{this.props.room.last_message.message && this.props.room.last_message.message}</p>
		    <p className="last-message-date">{this.props.room.last_message.sent && this.formatDate(this.props.room.last_message.sent)}</p>
		</div>
	    </div>
	);
    }
}

export default Room
