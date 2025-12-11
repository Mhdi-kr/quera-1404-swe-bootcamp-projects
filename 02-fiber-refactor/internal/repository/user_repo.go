package repository

import (
	"example.com/authorization/internal/repository/entity"
)

type UserRepository struct {
	Users []entity.User
}

func NewUserRepository() UserRepository {
	ur := UserRepository{}

	ur.Insert("mehdi", "1234")
	ur.Insert("amir", "1234")

	return ur
}

func (ur *UserRepository) Insert(username string, hashedPassword string) error {

	ur.Users = append(ur.Users, entity.User{
		Username: username,
		Password: hashedPassword,
	})

	return nil
}

func (ur *UserRepository) Search(username string) *entity.User {
	for _, u := range ur.Users {
		if u.Username == username {
			return &u
		}
	}

	return nil
}
