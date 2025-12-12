package pkg

import "database/sql"

type SQLRepository struct {
	DB *sql.DB
}

func NewSQLRepository() (SQLRepository, error) {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		return SQLRepository{}, err
	}

	return SQLRepository{
		DB: db,
	}, nil
}
