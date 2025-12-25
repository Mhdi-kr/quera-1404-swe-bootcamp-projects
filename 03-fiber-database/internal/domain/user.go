package domain

import "example.com/authorization/internal/repository/entity"

type User struct {
	Username string
	Email    string
}

func NewUserFromEntity(eu entity.User) User {
	return User{
		Username: eu.Username,
		Email:    eu.Email,
	}
}
