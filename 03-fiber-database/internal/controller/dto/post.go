package dto

import "time"

type Post struct {
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	URL         string     `json:"string"`
	Description string     `json:"description"`
}

type ProfilePostResponse struct {
	Posts []Post `json:"posts"`
}

type CreateProfilePostRequest struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}
