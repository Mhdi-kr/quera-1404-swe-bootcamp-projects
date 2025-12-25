package dto

import (
	"example.com/authorization/internal/domain"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewUserFromDomain(du domain.User) User {
	return User{
		Username: du.Username,
		Email:    du.Email,
	}
}
