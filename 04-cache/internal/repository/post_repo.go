package repository

import (
	"context"

	"example.com/authorization/internal/repository/entity"
	"example.com/authorization/pkg"
	"github.com/Masterminds/squirrel"
)

type PostRepository struct {
	Users   []entity.User
	sqlRepo pkg.SQLRepository
	cache   pkg.Cache
}

func NewPostRepository(sqlRepo pkg.SQLRepository, cache pkg.Cache) PostRepository {
	return PostRepository{
		sqlRepo: sqlRepo,
		Users:   []entity.User{},
		cache:   cache,
	}
}

func (ur *PostRepository) Insert(ctx context.Context, post entity.Post) (int64, error) {
	res, err := ur.sqlRepo.DB.ExecContext(ctx, "insert into `post` (`description`, `url`, `user_id`, `vote_count`) values (?, ?, ?, 0)", post.Description, post.URL, post.UserID)
	if err != nil {
		return 0, err
	}

	// we can add post to our search index here

	return res.LastInsertId()
}

func (ur *PostRepository) List(ctx context.Context, userID *int64, size uint64, page uint64) ([]entity.Post, error) {
	var posts []entity.Post

	query := squirrel.
		Select(
			"post.*",
			"COUNT(DISTINCT user_post_upvote.post_id) AS upvote_count",
			"COUNT(DISTINCT comment.id) AS comment_count",
		).
		From("post").
		LeftJoin("user_post_upvote ON post.id = user_post_upvote.post_id").
		LeftJoin("comment ON comment.post_id = post.id").
		GroupBy("post.id").
		Limit(size).
		Offset((page - 1) * size)

	if userID != nil {
		query = query.Where("post.user_id = ?", *userID)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return posts, err
	}

	rows, err := ur.sqlRepo.DB.QueryxContext(ctx, sql, args...)
	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var post entity.Post
		if err := rows.StructScan(&post); err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}

	return posts, nil

}

func (ur *PostRepository) DeleteByID(ctx context.Context, userID int64, postID int64) error {
	query := squirrel.Delete("post").Where(squirrel.And{
		squirrel.Eq{
			"id": postID,
		},
		squirrel.Eq{
			"user_id": userID,
		},
	})
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	result, err := ur.sqlRepo.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrPostNotFound
	}

	return nil
}
