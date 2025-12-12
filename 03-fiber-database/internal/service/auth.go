package service

import (
	"fmt"
	"time"

	"example.com/authorization/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type TokenString string

// we can add custom behaviour to our sub types
func (ts TokenString) Print() {
	fmt.Println(ts)
}

type AuthService struct {
	jwtSecret string
	userRepo  repository.UserRepository
}

func NewAuthorizationService(jwtSecret string, userRepo repository.UserRepository) AuthService {
	return AuthService{
		jwtSecret: jwtSecret,
		userRepo:  userRepo,
	}
}

func (as AuthService) GenerateToken(userId string) (TokenString, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		// A usual scenario is to set the expiration time relative to the current time
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "test",
		Subject:   userId,
		ID:        "0",
		Audience:  []string{},
	})

	tokenString, err := token.SignedString([]byte(as.jwtSecret))
	if err != nil {
		return "", err
	}

	return TokenString(tokenString), nil
}

func (as AuthService) ValidateToken(token string) (jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(as.jwtSecret), nil
	})
	if err != nil {
		return jwt.Token{}, err
	}

	return *jwtToken, nil
}

// TODO: move this to service layer
// hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
// if err != nil {
// 	return err
// }

// func checkUser(user entity.User) bool {
// 	isUserFound := false

// 	for _, u := range usersDB {
// 		resCompare := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
// 		fmt.Println(resCompare)
// 		if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)) == nil && user.Username == u.Username {
// 			isUserFound = true
// 		}
// 	}

// 	return isUserFound
// }
