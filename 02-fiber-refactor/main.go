package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserDetails struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type SelfResponseDTO struct {
	User   UserDetails `json:"user"`
	Status string      `json:"status"`
}

type LoginResponseDTO struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
	Token   string `json:"token"`
}

var usersDB []User

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

func searchUser(username string) *UserDetails {
	for _, u := range usersDB {
		if u.Username == username {
			return &UserDetails{
				Username: u.Username,
			}
		}
	}

	return nil
}

func readJwtSecret() []byte {
	return []byte("250703a7f4ec7712490ba2785b2538a71136d4c5200b2adc42bfd53225f2712a1bda5766")
}

func validateJwtToken(token string) (string, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return readJwtSecret(), nil
	})
	if err != nil {
		return "", err
	}

	userId, err := jwtToken.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	return userId, err
}

func createJwtToken(user User) (string, error) {
	sigingSecret := readJwtSecret()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		// A usual scenario is to set the expiration time relative to the current time
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "test",
		Subject:   user.Username,
		ID:        "0",
		Audience:  []string{},
	})
	tokenString, err := token.SignedString(sigingSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func main() {
	usersDB = append(usersDB, seedUser("mehdi", "1234"), seedUser("amir", "1234"))
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// r.Response.Write()
	})

	mux.HandleFunc("/api/v1/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		// handle serde (deser)
		user := User{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&user)

		usersDB = append(usersDB, seedUser(user.Username, user.Password))
	})

	mux.HandleFunc("/api/v1/self", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}

		authHeader := r.Header.Get("Authorization")

		if len(authHeader) == 0 {
			return
		}

		userId, err := validateJwtToken(authHeader)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
		}

		userDetails := searchUser(userId)
		if userDetails == nil {
			w.WriteHeader(http.StatusNotFound)
		}

		responseDto := SelfResponseDTO{
			User:   *userDetails,
			Status: "ok",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseDto)
	})

	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		responseDto := LoginResponseDTO{}
		user := User{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&user)

		isUserFound := checkUser(user)

		if isUserFound {
			jwtToken, err := createJwtToken(user)
			if err != nil {
				panic(err)
			}
			responseDto = LoginResponseDTO{
				Message: "ok",
				Error:   nil,
				Token:   jwtToken,
			}
		} else {
			responseDto = LoginResponseDTO{
				Message: "failed",
				Error:   errors.New("check credintials"),
				Token:   "",
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseDto)
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := srv.ListenAndServe()
	fmt.Println(err)
}
