CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    capacity INT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(255) PRIMARY KEY,
    password VARCHAR(255) NOT NULL,
    room_id INT REFERENCES rooms(id)
);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    message TEXT NOT NULL,
    author VARCHAR(255) REFERENCES users(username),
    date_sent TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    room_id INT REFERENCES rooms(id)
);

CREATE TABLE IF NOT EXISTS room_users (
    room_id INT REFERENCES rooms(id),
    username VARCHAR(255) REFERENCES users(username),
    PRIMARY KEY (room_id, username)
);
