package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func Db(cs *string, driverName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, *cs)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
