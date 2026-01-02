package entity

import (
	"database/sql"
)

type User struct {
	Id        int64          `db:"id"`
	Username  string         `db:"username"`
	Password  string         `db:"password"`
	Email     sql.NullString `db:"email"`
	FullName  sql.NullString `db:"full_name"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

type Users []User

func (u *User) IsValid() bool {
	return len(u.Password) > 0 && len(u.Username) > 0
}
