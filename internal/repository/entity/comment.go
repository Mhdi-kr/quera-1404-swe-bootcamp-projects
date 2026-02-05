package entity

import (
	"database/sql"
	"strconv"
	"time"
)

type Comment struct {
	Id        int64        `db:"id"`
	UserID    int64        `db:"user_id"`
	PostID    int64        `db:"post_id"`
	Content   string       `db:"content"`
	VoteCount uint64       `db:"vote_count"`
	CreatedAt sql.NullTime `db:"created_at" redis:"-"`
	UpdatedAt sql.NullTime `db:"updated_at" redis:"-"`
}

func (c Comment) ToHsetArgs() []string {
	return []string{
		strconv.FormatInt(c.Id, 10),
		strconv.FormatInt(c.UserID, 10),
		strconv.FormatInt(c.PostID, 10),
		c.Content,
		strconv.FormatUint(c.VoteCount, 10),
		c.CreatedAt.Time.Format(time.DateTime),
		c.UpdatedAt.Time.Format(time.DateTime),
	}
}

func NewCommentFromMap(in map[string]any) Comment {
	id := in["Id"].(int64)
	userId := in["UserID"].(int64)
	postId := in["PostID"].(int64)
	content := in["Content"].(string)
	voteCount := in["VoteCount"].(uint64)
	createdAt := in["CreatedAt"].(string)
	updatedAt := in["CreatedAt"].(string)
	ca, _ := time.Parse(time.DateTime, createdAt)
	ua, _ := time.Parse(time.DateTime, updatedAt)

	return Comment{
		Id:        id,
		UserID:    userId,
		PostID:    postId,
		Content:   content,
		VoteCount: voteCount,
		CreatedAt: sql.NullTime{
			Time:  ca,
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  ua,
			Valid: true,
		},
	}
}

type CommentUpvotesAgg struct {
	CommentId int64  `db:"commentId"`
	Upvotes   uint64 `db:"upvotes"`
}
