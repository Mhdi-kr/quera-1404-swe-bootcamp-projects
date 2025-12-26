package entity

import (
	"database/sql"
)

type Post struct {
	Id          int64        `db:"id"`
	Description string       `db:"description"`
	URL         string       `db:"url"`
	UserID      int64        `db:"user_id"`
	VoteCount   uint64       `db:"vote_count"`
	CreatedAt   sql.NullTime `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
}
