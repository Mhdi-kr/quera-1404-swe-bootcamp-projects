package pkg

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SQLRepository struct {
	DB *sqlx.DB
}

func NewSQLRepository(connectionUri string) (SQLRepository, error) {
	db, err := sql.Open("mysql", connectionUri)
	if err != nil {
		return SQLRepository{}, err
	}

	dbx := sqlx.NewDb(db, "mysql")

	return SQLRepository{
		DB: dbx,
	}, nil
}
