package query

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/apetropolsky/pmc_bot/library"
)

// InitDB creates structures and fill database by tracks
func InitDB(db *sql.DB, rootpath string) {
	_, err := db.Exec(`DROP TABLE IF EXISTS tracks;`)
	if err != nil {
		fmt.Println("Не удалил таблицу")
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tracks (
		name varchar,
		path varchar UNIQUE,
		category varchar,
		artist varchar,
		album varchar
		);`)
	if err != nil {
		fmt.Println("Не смог создать таблицу")
	}

	tracks, _ := library.GetContent(db, rootpath)
	for _, track := range tracks {
		insertTrack(db, track)
	}
}

// insertTrack put given track info into database
func insertTrack(db *sql.DB, track library.Track) {
	insStr := fmt.Sprintf(
		`INSERT INTO tracks VALUES 
		($$%[1]s$$, $$%[2]s$$, $$%[3]s$$, $$%[4]s$$, $$%[5]s$$) 
		ON CONFLICT (path) DO UPDATE 
		SET name = $$%[1]s$$, category = $$%[3]s$$, artist = $$%[4]s$$, album = $$%[5]s$$;`,
		track.Name,
		track.Path,
		track.Category,
		track.Artist,
		track.Album,
	)

	_, err := db.Exec(insStr)
	if err != nil {
		panic(err)
	}
}

// Category provide category contents
func Category(db *sql.DB, category string) {
	category = "'%" + category + "%'"

	request := fmt.Sprintf(
		`SELECT DISTINCT artist
			FROM tracks WHERE 
			LOWER(category) LIKE LOWER(%s);`,
		category,
	)

	ContentQuery(db, request)
}

// Artist provide artist albums
func Artist(db *sql.DB, artist string) {
	artist = "'%" + artist + "%'"

	request := fmt.Sprintf(
		`SELECT DISTINCT album
			FROM tracks WHERE 
			LOWER(artist) LIKE LOWER(%s);`,
		artist,
	)

	ContentQuery(db, request)
}

// Album provide album tracks
func Album(db *sql.DB, album string) {
	album = "'%" + album + "%'"

	request := fmt.Sprintf(
		`SELECT name
			FROM tracks WHERE 
			LOWER(album) LIKE LOWER(%s);`,
		album,
	)

	ContentQuery(db, request)
}

// ContentQuery filter records by one column
func ContentQuery(db *sql.DB, reqStr string) {
	var column string
	var result []string

	rows, err := db.Query(reqStr)
	if err != nil {
		fmt.Println("Не смог обработать запрос")
		return
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&column)
		if err != nil {
			panic(err)
		}
		result = append(result, column)
	}

	sort.Strings(result)
	for _, rec := range result {
		fmt.Println(rec)
	}

}

// SimpleQuery get all records matches pattern and
// print trackname | artist | album
func SimpleQuery(db *sql.DB, toFind string) {
	var name string
	var artist string
	var album string
	var result []string

	toFind = "'%" + toFind + "%'"

	resSrt := fmt.Sprintf(
		`SELECT name, artist, album 
			FROM tracks WHERE 
			LOWER(name) LIKE LOWER(%[1]s) 
			OR LOWER(artist) LIKE LOWER(%[1]s)
			OR LOWER(album) LIKE LOWER(%[1]s);`,
		toFind,
	)

	rows, err := db.Query(resSrt)
	if err != nil {
		fmt.Println("Не смог обработать запрос")
		return
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name, &artist, &album)
		if err != nil {
			panic(err)
		}
		info := fmt.Sprintf("%s | %s | %s", name, artist, album)
		result = append(result, info)
	}

	sort.Strings(result)
	for _, rec := range result {
		fmt.Println(rec)
	}

}
