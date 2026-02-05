package dto

import "time"

type ListPostsRequest struct {
	Page uint64 `query:"page"`
	Size uint64 `query:"size"`
}

func (lpr *ListPostsRequest) Sanitize() {
	if lpr.Page == 0 {
		lpr.Page = 1
	}

	if lpr.Size > 100 {
		lpr.Page = 100
	}

	if lpr.Size == 0 {
		lpr.Size = 4
	}
}

type ListPostsResponse struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	Id               int        `json:"id"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt"`
	URL              string     `json:"url"`
	Description      string     `json:"description"`
	NumberOfComments uint64     `json:"numberOfComments"`
	NumberOfUpvotes  uint64     `json:"numberOfUpvotes"`
}

type ProfilePostResponse struct {
	Posts []Post `json:"posts"`
}

type CreateProfilePostRequest struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}
