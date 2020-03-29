package main

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/apetropolsky/pmc_bot/db"
	"github.com/apetropolsky/pmc_bot/structs"
	_ "github.com/lib/pq"
)

func initDB(db *sql.DB, rootpath string) {
	folders, _ := structs.GetContent(rootpath)

	for _, folder := range folders {
		name := folder.Name()
		path := fmt.Sprintf("%s/%s", rootpath, folder.Name())

		category := structs.Category{Name: name, Path: path}
		category.GetArtist(db)
	}
}

func main() {

	db, err := db.ConnectDB()
	initDB(db, "/media/admn/data/mp3")

	var artist string
	var name string
	var result []string

	resSrt := fmt.Sprintf(
		`SELECT DISTINCT artist, name 
		FROM tracks WHERE 
		category = $$%[1]s$$ AND album = $$%[1]s$$;`,
		"Rock",
	)
	rows, err := db.Query(resSrt)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&artist, &name)
		if err != nil {
			panic(err)
		}
		aa := artist + " - " + name
		result = append(result, aa)
	}

	sort.Strings(result)
	for _, rec := range result {
		fmt.Println(rec)
	}

	defer db.Close()
}
