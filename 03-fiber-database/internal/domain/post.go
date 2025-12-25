package domain

import (
	"time"

	"example.com/authorization/internal/repository/entity"
)

type PostFilters struct {
	Page uint64
	Size uint64
}

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

type Posts []Post

func NewPostsFromEntities(p []entity.Post) []Post {
	var posts []Post
	for _, pe := range p {
		posts = append(posts, Post{
			Id:          pe.Id,
			Description: pe.Description,
			URL:         pe.URL,
			UserID:      pe.UserID,
			VoteCount:   pe.VoteCount,
			CreatedAt:   pe.CreatedAt.Time,
			UpdatedAt:   pe.UpdatedAt.Time,
		})
	}

	return posts
}
