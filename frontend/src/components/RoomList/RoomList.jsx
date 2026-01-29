import React, { Component } from "react";
import Room from '../Room/Room';
import './RoomList.scss';

class RoomList extends Component {
    constructor(props) {
	super(props);
    }

    render() {
	return (
	    <div id="room-list">
		<h2>Rooms</h2>

		<div id="chatroom-selection">
		{[...this.props.rooms.values()].map((room) => (
		    <Room key={room.id} 
			room={room}
			changeRoom={this.props.changeRoom}
		    />)
		)}
		</div>
	    </div>
	);
    }
}

export default RoomList;
