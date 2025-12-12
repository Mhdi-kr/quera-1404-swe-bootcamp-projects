package dto

import "time"

type Post struct {
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
}

type ProfilePostResponse struct {
	Posts []Post `json:"posts"`
}
