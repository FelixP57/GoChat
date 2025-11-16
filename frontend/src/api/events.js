
export class Event {
    constructor(type, payload) {
	this.type = type;
	this.payload = payload;
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
    constructor(id, name, last_message) {
	this.id = id;
	this.name = name;
	this.last_message = last_message;
    }
}

export class CreateRoomEvent {
    constructor(username) {
	this.username = username
    }
}

