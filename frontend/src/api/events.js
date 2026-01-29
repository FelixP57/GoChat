
export class Event {
    constructor(type, payload) {
	this.type = type;
	this.payload = payload;
    }
}

export class RoomUser {
    constructor(username, online, typing) {
	this.username = username;
	this.online = online;
	this.typing = typing;
    }
}

export class SendMessageEvent {
    constructor(message, room_id) {
	this.message = message;
	this.room_id = room_id;
    }
}

export class NewMessageEvent {
    constructor(message, from, sent, room_id) {
	this.message = message;
	this.from = from;
	this.sent = sent;
	this.room_id = room_id;
    }
}

export class NewRoomEvent {
    constructor(id, name, last_message, users) {
	this.id = id;
	this.name = name;
	this.users = users
	this.last_message = last_message;
    }
}

export class CreateRoomEvent {
    constructor(username) {
	this.username = username
    }
}

export class UserConnectedEvent {
    constructor(username) {
	this.username = username
    }
}

export class UserDisconnectedEvent {
    constructor(username) {
	this.username = username
    }
}
