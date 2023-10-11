package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func Db(cs *string) (*sql.DB, error) {
	db, err := sql.Open("postgres", *cs)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
