package dto

import "time"

type RegisterRequest struct {
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	Email     string     `json:"email"`
	FullName  *string    `json:"firstName"`
	BrithDate *time.Time `json:"birthDate"`
}
