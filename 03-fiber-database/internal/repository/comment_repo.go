package repository

import (
	"context"

	"example.com/authorization/internal/repository/entity"
	"example.com/authorization/pkg"
	"github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
)

const MYSQL_KEY_EXITS uint16 = 1062

type CommentRepo struct {
	Users   []entity.User
	sqlRepo pkg.SQLRepository
}

func NewCommentRepo(sqlRepo pkg.SQLRepository) CommentRepo {
	return CommentRepo{
		sqlRepo: sqlRepo,
		Users:   []entity.User{},
	}
}

func (ur *CommentRepo) Insert(ctx context.Context, comment entity.Comment) (int64, error) {
	sql, args, err := squirrel.Insert("comment").Columns(
		"user_id",
		"post_id",
		"content",
		"vote_count",
	).Values(
		comment.UserID,
		comment.PostID,
		comment.Content,
		0,
	).ToSql()
	if err != nil {
		return 0, err
	}

	res, err := ur.sqlRepo.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (ur *CommentRepo) List(ctx context.Context, postID int64, size uint64, page uint64) ([]entity.Comment, error) {
	var comments []entity.Comment
	sql, args, err := squirrel.Select(
		"comment.id as id",
		"count(user_comment_upvote.comment_id) as vote_count",
		"comment.user_id",
		"comment.post_id",
		"comment.content",
	).
		From("comment").
		Limit(size).
		Join("user_comment_upvote on comment.id = user_comment_upvote.comment_id").
		GroupBy("comment.id").
		Offset((page-1)*size).
		Where("post_id = ?", postID).
		ToSql()

	rows, err := ur.sqlRepo.DB.QueryxContext(ctx, sql, args...)
	if err != nil {
		return comments, err
	}

	for rows.Next() {
		var comment entity.Comment
		err := rows.StructScan(&comment)
		if err != nil {
			return comments, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (ur *CommentRepo) DeleteByID(ctx context.Context, userID int64, commentID int64) error {
	query := squirrel.Delete("comment").Where(squirrel.And{
		squirrel.Eq{
			"id": commentID,
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
		return ErrCommentNotFound
	}

	return nil
}

func (ur *CommentRepo) Upvote(ctx context.Context, userID int64, commentID int64) (bool, error) {
	state := false
	sqlstr, args, err := squirrel.Insert("user_comment_upvote").Columns("user_id", "comment_id").Values(
		userID,
		commentID,
	).ToSql()
	if err != nil {
		return state, err
	}

	_, err = ur.sqlRepo.DB.ExecContext(ctx, sqlstr, args...)
	if err == nil {
		state = true
		return state, nil
	}

	mysqlerr, ok := err.(*mysql.MySQLError)
	if !ok {
		return state, err
	}

	if mysqlerr.Number == MYSQL_KEY_EXITS {
		delsqlstr, delargs, delerr := squirrel.Delete("user_comment_upvote").Where("user_id = ?", userID).Where("comment_id = ?", commentID).ToSql()
		_, delerr = ur.sqlRepo.DB.ExecContext(ctx, delsqlstr, delargs...)
		if delerr != nil {
			return state, delerr
		}
	}

	return state, nil
}

// func (ur *CommentRepo) Upvote(ctx context.Context, commentID int64) error {
// 	tx, err := ur.sqlRepo.DB.BeginTxx(ctx, &sql.TxOptions{})
// 	ssql, sargs, err := squirrel.
// 		Select("*").
// 		From("comment").
// 		Where("id = ?", commentID).
// 		ToSql()
// 	if err != nil {
// 		return err
// 	}

// 	rows, err := tx.QueryxContext(ctx, ssql, sargs...)
// 	if err != nil {
// 		return err
// 	}

// 	var comments []entity.Comment
// 	for rows.Next() {
// 		var c entity.Comment
// 		err := rows.StructScan(&c)
// 		comments = append(comments, c)
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}
// 	}

// 	if len(comments) == 0 {
// 		return ErrCommentNotFound
// 	}

// 	usql, uargs, err := squirrel.
// 		Update("comment").
// 		Set("vote_count", comments[0].VoteCount+1).
// 		Where("id = ?", commentID).
// 		ToSql()
// 	if err != nil {
// 		return errors.Join(err, tx.Rollback())
// 	}

// 	result, err := tx.ExecContext(ctx, usql, uargs...)
// 	if err != nil {
// 		return errors.Join(err, tx.Rollback())
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return errors.Join(err, tx.Rollback())
// 	}

// 	if rowsAffected == 0 {
// 		return errors.Join(ErrCommentNotFound, tx.Rollback())
// 	}

// 	return tx.Commit()
// }
