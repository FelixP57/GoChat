import React, { Component } from "react";
import './RoomCreationForm.scss';

class RoomCreationForm extends Component {
    constructor(props) {
	super(props);
	this.state = {
	    username: '',
	};
	this.createRoom = this.createRoom.bind(this);
    }

    createRoom(event) {
	event.preventDefault();
	if (this.state.username != '') {
	    this.props.createRoom(this.state.username);
	}
	return false;
    }

    render() {
	return (
	    <div id="room-creation">
		<form onSubmit={this.createRoom} >
		    <label htmlFor="with-username">New room with:</label>
		    <br />
		    <input type="text" id="with-username" name="with-username" onChange={(e) => this.setState({username: e.target.value})} value={this.state.username} /><br />
		    <input className="submit" type="submit" value="Create room" />
		</form>
	    </div>
	);
    }
}

export default RoomCreationForm;
