package repository

import (
	"context"
	"fmt"

	"example.com/authorization/internal/repository/entity"
	"example.com/authorization/pkg"
)

type UserRepository struct {
	Users   []entity.User
	sqlRepo pkg.SQLRepository
}

func NewUserRepository(sqlRepo pkg.SQLRepository) UserRepository {
	return UserRepository{
		sqlRepo: sqlRepo,
		Users:   []entity.User{},
	}
}

func (ur *UserRepository) Insert(ctx context.Context, username string, hashedPassword string) error {
	query := fmt.Sprintf("insert into `user` (`password`, `username`) values ('%s', '%s')", hashedPassword, username)
	_, err := ur.sqlRepo.DB.ExecContext(ctx, query)

	return err
}

func (ur *UserRepository) GetOneByUsername(ctx context.Context, username string) (entity.User, error) {
	var users []entity.User

	err := ur.sqlRepo.DB.SelectContext(ctx, &users, "select * from user where username = ?", username)
	if err != nil {
		return entity.User{}, err
	}

	if len(users) == 0 {
		return entity.User{}, ErrUserNotFound
	}

	return users[0], nil
}

func (ur *UserRepository) ListAll(ctx context.Context) ([]entity.User, error) {
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
