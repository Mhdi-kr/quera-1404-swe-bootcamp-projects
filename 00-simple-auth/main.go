package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Password string
	Username string
}

func seedUser(username string, password string) User {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	// REFACTOR: handle error
	if err != nil {
		panic(err)
	}

	return User{
		Username: username,
		Password: string(hashBytes),
	}
}

func checkUser(user User) bool {
	var usersDB []User
	usersDB = append(usersDB, seedUser("mehdi", "1234"), seedUser("amir", "1234"))

	isUserFound := false

	for _, u := range usersDB {
		resCompare := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
		fmt.Println(resCompare)
		if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)) == nil && user.Username == u.Username {
			isUserFound = true
		}
	}

	return isUserFound
}

func main() {
	// cmdReader := bufio.NewReader(os.Stdin)
	fmt.Println("enter username:")
	// inputUsername, _ := cmdReader.ReadString('\n')
	fmt.Println("enter password:")
	// inputPassword, _ := cmdReader.ReadString('\n')

	inputUser := User{
		Username: "mehdi",
		Password: "12345",
	}

	fmt.Println("input user: ", inputUser)

	isUserFound := checkUser(inputUser)

	if isUserFound {
		fmt.Println("hello", inputUser.Username)
	} else {
		fmt.Println("check credentials", inputUser.Username)
	}
	// no layering
	// we want to create a user manager without database
	// authorization
}
