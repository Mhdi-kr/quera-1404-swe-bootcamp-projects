package entity

import (
	"database/sql"
)

type Post struct {
	Id           int64        `db:"id"`
	Description  string       `db:"description"`
	URL          string       `db:"url"`
	UserID       int64        `db:"user_id"`
	UpvoteCount  uint64       `db:"upvote_count"`
	CommentCount uint64       `db:"comment_count"`
	CreatedAt    sql.NullTime `db:"created_at"`
	UpdatedAt    sql.NullTime `db:"updated_at"`
}
