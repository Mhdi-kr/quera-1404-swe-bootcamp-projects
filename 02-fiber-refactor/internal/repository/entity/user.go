package entity

type User struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func (u *User) IsEmpty() bool {
	return len(u.Password) > 0 && len(u.Username) > 0
}
