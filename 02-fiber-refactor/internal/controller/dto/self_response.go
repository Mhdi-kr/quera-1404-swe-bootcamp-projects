package dto

type SelfResponse struct {
	Details UserDetails `json:"details"`
	Status  string      `json:"status"`
}
