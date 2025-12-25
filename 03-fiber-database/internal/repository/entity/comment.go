package entity

import "database/sql"

type Comment struct {
	Id        int64        `db:"id"`
	UserID    int64        `db:"user_id"`
	PostID    int64        `db:"post_id"`
	Content   string       `db:"content"`
	VoteCount uint64       `db:"vote_count"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
