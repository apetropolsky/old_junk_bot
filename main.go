package main

import (
	"database/sql"
	"os"

	"github.com/apetropolsky/pmc_bot/query"
	_ "github.com/lib/pq"
)

// connectDB is a simple database connector
func connectDB() (*sql.DB, error) {
	connStr := "dbname=test user=postgres host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	return db, err
}

func main() {
	db, _ := connectDB()
	query.InitDB(db, "/media/admn/data/mp3")

	query.Album(db, os.Args[1])
	defer db.Close()
}
