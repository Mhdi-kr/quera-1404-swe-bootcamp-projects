package dto

import "time"

type ListCommentsRequest struct {
	Page uint64 `query:"page"`
	Size uint64 `query:"size"`
}

func (lcr *ListCommentsRequest) Sanitize() {
	if lcr.Page == 0 {
		lcr.Page = 1
	}

	if lcr.Size > 100 {
		lcr.Page = 100
	}

	if lcr.Size == 0 {
		lcr.Size = 4
	}
}

type ListCommentsResponse struct {
	Comments []Comment `json:"comments"`
}

type Comment struct {
	Id        int        `json:"id"`
	UserID    int64      `json:"userId"`
	PostID    int64      `json:"postId"`
	Content   string     `json:"content"`
	VoteCount uint64     `json:"voteCount"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type CreateCommentRequest struct {
	Content string `json:"content"`
}
