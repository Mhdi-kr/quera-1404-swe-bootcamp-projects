package repository

import (
	"context"
	"fmt"

	"example.com/authorization/internal/repository/entity"
	"example.com/authorization/pkg"
)

type PostRepository struct {
	Users   []entity.User
	sqlRepo pkg.SQLRepository
}

func NewPostRepository(sqlRepo pkg.SQLRepository) PostRepository {
	return PostRepository{
		sqlRepo: sqlRepo,
		Users:   []entity.User{},
	}
}

func (ur *PostRepository) Insert(username string, hashedPassword string) error {
	ur.Users = append(ur.Users, entity.User{
		Username: username,
		Password: hashedPassword,
	})

	return nil
}

func (ur *PostRepository) GetOneByUsername(username string) (entity.User, error) {
	for _, u := range ur.Users {
		if u.Username == username {
			return u, nil
		}
	}

	return entity.User{}, ErrUserNotFound
}

func (ur *PostRepository) ListAll(ctx context.Context) ([]entity.User, error) {
	rows, err := ur.sqlRepo.DB.QueryxContext(ctx, "select * from user")
	if err != nil {
		return make([]entity.User, 0), err
	}

	var users entity.Users
	for rows.Next() {
		var usr entity.User
		err := rows.StructScan(&usr)
		if err != nil {
			fmt.Println(err)
		}
		users = append(users, usr)
	}

	return users, nil
}
