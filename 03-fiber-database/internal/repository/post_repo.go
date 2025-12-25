package repository

import (
	"context"

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

func (ur *PostRepository) Insert(ctx context.Context, post entity.Post) (int64, error) {
	res, err := ur.sqlRepo.DB.ExecContext(ctx, "insert into `post` (`description`, `url`, `user_id`, `vote_count`) values (?, ?, ?, 0)", post.Description, post.URL, post.UserID)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}
