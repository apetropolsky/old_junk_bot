package db

import (
	"database/sql"
)

// ConnectDB is a simple database connector
func ConnectDB() (*sql.DB, error) {
	connStr := "dbname=test user=postgres host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	return db, err
}
