package main

import "fmt"

type User struct {
	Password string
	Username string
}

var usersDB []User

func main() {
	usersDB = append(usersDB, User{
		Password: "1234",
		Username: "mehdi",
	}, User{
		Password: "1234",
		Username: "amir",
	})

	candidateUser := User{
		Password: "1234",
		Username: "fateme",
	}

	isUserFound := false
	foundUser := User{}
	for _, user := range usersDB {
		if candidateUser.Password == user.Password && candidateUser.Username == user.Username {
			foundUser = user
			isUserFound = true
		}
	}
	if isUserFound {
		fmt.Println("hello", foundUser.Username)
	} else {
		fmt.Println("check credentials", foundUser.Username)
	}
	// no layering
	// we want to create a user manager without database
	// authorization
}
