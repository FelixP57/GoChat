import React, { Component } from "react";
import './Profile.scss';

class Profile extends Component {
    constructor(props) {
	super(props);
    }

    render() {
	return (
	    <div id="profile">
		<h5 id="profile-username">{this.props.username}</h5>
		<form id="logout" onSubmit={this.props.logout}>
		    <button type="submit">Logout</button>
		</form>
	    </div>
	);
    }
}

export default Profile;

