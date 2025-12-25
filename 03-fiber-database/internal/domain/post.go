package domain

import (
	"time"

	"example.com/authorization/internal/repository/entity"
)

type Post struct {
	Id          int64
	Description string
	URL         string
	UserID      int64
	VoteCount   uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewPostFromEntity(p entity.Post) Post {
	return Post{
		Id:          p.Id,
		Description: p.Description,
		URL:         p.URL,
		UserID:      p.UserID,
		VoteCount:   p.VoteCount,
		CreatedAt:   p.CreatedAt.Time,
		UpdatedAt:   p.UpdatedAt.Time,
	}
}
