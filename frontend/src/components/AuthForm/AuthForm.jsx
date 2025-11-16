import React, { Component } from "react";
import "./AuthForm.scss";

class AuthForm extends Component {
    constructor(props) {
	super(props);
	this.loginHandler = props.loginHandler;
	this.signupHandler = props.signupHandler;
    }

    render() {
	return (
	    <div id="auth">
		<h3>Authentification</h3>
		<form id="auth-form">
		    <label htmlFor="username">Username:</label>
		    <br />
		    <input type="text" id="username" name="username" /><br />
		    <br />
		    <label htmlFor="password">Password:</label>
		    <br />
		    <input type="password" id="password" name="password" /><br /><br />
		</form>
		<div id="auth-submit">
		    <form onSubmit={this.loginHandler} id="login-form">
			<input className="submit" type="submit" value="Login" />
		    </form>
		    <form onSubmit={this.signupHandler} id="signup-form">
			<input className="submit" type="submit" value="Sign up" />
		    </form>
		</div>
	    </div>
	);
    }
}

export default AuthForm;
