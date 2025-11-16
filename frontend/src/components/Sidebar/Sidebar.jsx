import React, { Component } from "react";
import RoomCreationForm from '../RoomCreationForm/RoomCreationForm';
import RoomList from '../RoomList/RoomList';
import Profile from '../Profile/Profile';
import './Sidebar.scss';

class Sidebar extends Component {
    constructor(props) {
	super(props);
    }

    render() {
	return (
	    <div className="sidebar">
		<RoomList rooms={this.props.rooms} selectedRoom={this.props.selectedRoom} changeRoom={this.props.changeRoom} />
		
		<br />
		<br />

		<RoomCreationForm createRoom={this.props.createRoom} />

		<Profile username={this.props.username} logout={this.props.disconnect} />
	    </div>
	);
    }
}

export default Sidebar;
