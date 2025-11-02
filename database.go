package main

import (
	"database/sql"
	"os"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	dbname = "gochat_db"
)

type Database struct {
	db *sql.DB
	hub *Hub
}

func (db *Database) closeDb() {
	db.db.Close()
}

func getDb(hub *Hub) *Database {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, os.Getenv("PSQL_PWD"), dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return &Database{db: db, hub: hub}
}

func (db *Database) addUser(user *User) {
	sqlStatement := `INSERT INTO users (username, password, room_id) VALUES ($1, $2, $3);`
	_, err := db.db.Exec(sqlStatement, user.username, user.password, user.roomId)
	if err != nil {
		panic(err)
	}
}

func (db *Database) getUserByUsername(username string) (*User, error) {
	sqlStatement := `SELECT * FROM users WHERE username=$1;`
	var user User
	row := db.db.QueryRow(sqlStatement, username)
	err := row.Scan(&user.username, &user.password, &user.roomId)
	switch err {
	case sql.ErrNoRows:
		return nil, err
	case nil:
		return &user, nil
	default:
		panic(err)
	}
}

func (db *Database) updateUserRoom(user *User, roomId int) {
	sqlStatement := `UPDATE users SET room_id=$1 WHERE username=$2;`
	_, err := db.db.Exec(sqlStatement, roomId, user.username);
	if err != nil {
		panic(err)
	}
}

func (db *Database) addRoom(room *Room) int {
	sqlStatement := `INSERT INTO rooms (capacity, name) VALUES ($1, $2) RETURNING Id;`
	var id int
	row := db.db.QueryRow(sqlStatement, room.capacity, room.name)
	err := row.Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}

func (db *Database) getRooms() (map[int]*Room, error) {
	sqlStatement := `SELECT id, capacity, name FROM rooms;`
	rooms := make(map[int]*Room)
	rows, err := db.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		room := newRoom(db.hub)
		var id int
		err = rows.Scan(&id, &room.capacity, &room.name)
		if err != nil {
			return nil, err	
		}
		room.users = db.getRoomUsers(id)
		rooms[id] = room
	}
	return rooms, nil
}

func (db *Database) updateRoomName(roomId int, name string) {
	sqlStatement := `UPDATE rooms SET name=$1 WHERE id=$2;`
	_, err := db.db.Exec(sqlStatement, name, roomId)
	if err != nil {
		panic(err)
	}
}

func (db *Database) addMessage(message NewMessageEvent, roomId int) {
	sqlStatement := `INSERT INTO messages (message, author, date_sent, room_id) VALUES ($1, $2, $3, $4);`
	_, err := db.db.Exec(sqlStatement, message.Message, message.From, message.Sent, roomId)
	if err != nil {
		panic(err)
	}
}

func (db *Database) getMessages(roomId int) []NewMessageEvent {
	sqlStatement := `SELECT * FROM messages WHERE room_id=$1;`
	var events []NewMessageEvent
	rows, err := db.db.Query(sqlStatement, roomId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var event NewMessageEvent
		var id int
		err = rows.Scan(&event.Message, &event.From, &event.Sent, &id)
		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}
	return events
}

func (db *Database) getRoomUsers(roomId int) map[*User]bool {
	sqlStatement := `SELECT username FROM room_users WHERE room_id=$1;`
	users := make(map[*User]bool) 
	rows, err := db.db.Query(sqlStatement, roomId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		var user *User
		err = rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		user, err = db.getUserByUsername(username)
		if err != nil {
			panic(err)
		}
		users[user] = true
	}
	return users
}

func (db *Database) addUserToRoom(username string, roomId int) {
	sqlStatement := `INSERT INTO room_users (room_id, username) VALUES ($1, $2);`
	_, err := db.db.Exec(sqlStatement, roomId, username)
	if err != nil {
		panic(err)
	}
}

