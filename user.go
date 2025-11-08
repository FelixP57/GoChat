package main


type User struct {
	username string

	password string

}

func newUser(username string, password string) *User {
	return &User{
		username: username,
		password: password,
	}
}

