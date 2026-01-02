package domain

import (
	"database/sql"
	"time"

	"example.com/authorization/internal/controller/dto"
	"example.com/authorization/internal/repository/entity"
)

type CommentFilters struct {
	Page uint64
	Size uint64
}

type Comment struct {
	Id        int64
	UserID    int64
	PostID    int64
	Content   string
	VoteCount uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Comment) ToEntity() entity.Comment {
	return entity.Comment{
		Id:        c.Id,
		UserID:    c.UserID,
		PostID:    c.PostID,
		Content:   c.Content,
		VoteCount: c.VoteCount,
		CreatedAt: sql.NullTime{Time: c.CreatedAt, Valid: true},
		UpdatedAt: sql.NullTime{Time: c.UpdatedAt, Valid: true},
	}
}

func (c *Comment) ToDTO() dto.Comment {
	return dto.Comment{
		Id:        int(c.Id),
		UserID:    c.UserID,
		PostID:    c.PostID,
		Content:   c.Content,
		VoteCount: c.VoteCount,
		CreatedAt: c.CreatedAt,
		UpdatedAt: &c.UpdatedAt,
	}
}

func NewCommentFromEntity(p entity.Comment) Comment {
	return Comment{
		Id:        p.Id,
		PostID:    p.PostID,
		Content:   p.Content,
		UserID:    p.UserID,
		VoteCount: p.VoteCount,
		CreatedAt: p.CreatedAt.Time,
		UpdatedAt: p.UpdatedAt.Time,
	}
}

func NewCommentsFromEntities(ces []entity.Comment) []Comment {
	var comments []Comment
	for _, ce := range ces {
		comments = append(comments, Comment{
			Id:        ce.Id,
			PostID:    ce.PostID,
			Content:   ce.Content,
			UserID:    ce.UserID,
			VoteCount: ce.VoteCount,
			CreatedAt: ce.CreatedAt.Time,
			UpdatedAt: ce.UpdatedAt.Time,
		})
	}

	return comments
}
