package repository

import (
	"fmt"

	"example.com/authorization/internal/repository/entity"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	Users []entity.User
}

func NewUserRepository() UserRepository {
	ur := UserRepository{}

	hash1, _ := bcrypt.GenerateFromPassword([]byte("1234"), 4)
	hash2, _ := bcrypt.GenerateFromPassword([]byte("1234"), 4)
	ur.Insert("mehdi", string(hash1))
	ur.Insert("amir", string(hash2))

	return ur
}

func (ur *UserRepository) Insert(username string, hashedPassword string) error {
	ur.Users = append(ur.Users, entity.User{
		Username: username,
		Password: hashedPassword,
	})

	return nil
}

func (ur *UserRepository) GetOneByID(username string) (entity.User, error) {
	for _, u := range ur.Users {
		if u.Username == username {
			return u, nil
		}
	}

	return entity.User{}, ErrUserNotFound
}

func (ur *UserRepository) ListAll() ([]entity.User, error) {
	fmt.Println("listall")
	return ur.Users, nil
}
