package main


type User struct {
	username string

	password string

	roomId int 
}

func newUser(username string, password string) *User {
	return &User{
		username: username,
		password: password,
		roomId: 0,
	}
}

