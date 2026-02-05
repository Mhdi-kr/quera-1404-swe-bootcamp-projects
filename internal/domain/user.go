package domain

import (
	"example.com/authorization/internal/controller/dto"
	"example.com/authorization/internal/repository/entity"
)

type User struct {
	Username string
	Email    string
}

func (u *User) ToDTO() dto.User {
	return dto.User{
		Username: u.Username,
		Email:    u.Email,
	}
}

func NewUserFromEntity(eu entity.User) User {
	return User{
		Username: eu.Username,
		Email:    eu.Email.String,
	}
}
