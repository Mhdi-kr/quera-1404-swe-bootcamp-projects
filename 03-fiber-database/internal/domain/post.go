package domain

import (
	"database/sql"
	"time"

	"example.com/authorization/internal/controller/dto"
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

func (p *Post) ToEntity() entity.Post {
	return entity.Post{
		Id:          p.Id,
		Description: p.Description,
		URL:         p.URL,
		UserID:      p.UserID,
		VoteCount:   p.VoteCount,
		CreatedAt:   sql.NullTime{Time: p.CreatedAt, Valid: true},
		UpdatedAt:   sql.NullTime{Time: p.UpdatedAt, Valid: true},
	}
}

func (p *Post) ToDTO() dto.Post {
	return dto.Post{
		Id:          int(p.Id),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   &p.UpdatedAt,
		URL:         p.URL,
		Description: p.Description,
	}
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
