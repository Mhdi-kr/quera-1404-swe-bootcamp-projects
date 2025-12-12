package dto

type SelfResponse struct {
	User   User   `json:"user"`
	Status string `json:"status"`
}
