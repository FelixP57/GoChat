package main


type User struct {
	username string

	online bool
}

func newUser(username string) *User {
	return &User{
		username: username,
		online: false,
	}
}

