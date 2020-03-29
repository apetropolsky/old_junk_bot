package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	_ "github.com/lib/pq"
)

// Category on the top level of library
type Category struct {
	ID         int
	Name, Path string
	Artists    []Artist
}

func (c *Category) listContent() {
	for _, artist := range c.Artists {
		fmt.Println(artist.Name)
	}
}

func (c *Category) getArtist(db *sql.DB) {
	folders, _ := getContent(c.Path)

	for _, folder := range folders {
		name := folder.Name()
		path := fmt.Sprintf("%s/%s", c.Path, name)
		artist := Artist{Name: name, Path: path, Category: c.Name}
		artist.getAlbums(db)
		c.Artists = append(c.Artists, artist)
	}
}

// Artist in category
type Artist struct {
	Name, Path, Category string
	Albums               []Album
	Tracks               []Track
}

func (a *Artist) listContent() {
	for _, album := range a.Albums {
		fmt.Println(album.Name)
	}
}

func (a *Artist) getAlbums(db *sql.DB) {
	a.Tracks = getTracks(db, a.Path, a.Category, a.Category, a.Name)
	albums, _ := getContent(a.Path)

	for _, alb := range albums {
		name := alb.Name()
		path := fmt.Sprintf("%s/%s", a.Path, alb.Name())
		album := Album{Name: name, Path: path, Category: a.Category, Artist: a.Name}
		album.Tracks = getTracks(db, album.Path, a.Category, album.Name, a.Name)
		a.Albums = append(a.Albums, album)
	}
}

// Album in artist folder
type Album struct {
	Name, Path, Category, Artist string
	Tracks                       []Track
}

func (a *Album) listContent() {
	for _, track := range a.Tracks {
		fmt.Println(track.Name)
	}
}

// Track structure
type Track struct {
	Name, Path, Category, Album, Artist string
}

func getTracks(db *sql.DB, rootpath string, category string, album string, artist string) []Track {
	var tracks []Track
	_, files := getContent(rootpath)

	for _, file := range files {
		path := fmt.Sprintf("%s/%s", rootpath, file)
		track := Track{Name: file, Path: path, Category: category, Album: album, Artist: artist}
		tracks = append(tracks, track)
		insertTrack(db, track)
	}
	return tracks
}

func insertTrack(db *sql.DB, track Track) {
	insStr := fmt.Sprintf("INSERT INTO tracks VALUES ($$%[1]s$$, $$%[2]s$$, $$%[3]s$$, $$%[4]s$$, $$%[5]s$$) ON CONFLICT (path) DO UPDATE SET name = $$%[1]s$$, category = $$%[3]s$$, artist = $$%[4]s$$, album = $$%[5]s$$;", track.Name, track.Path, track.Category, track.Artist, track.Album)
	_, err := db.Exec(insStr)
	if err != nil {
		panic(err)
	}

}

func getContent(rootpath string) ([]os.FileInfo, []string) {
	var folders []os.FileInfo
	var files []string
	content, _ := ioutil.ReadDir(rootpath)
	for _, item := range content {
		if item.IsDir() {
			folders = append(folders, item)
		} else {
			files = append(files, item.Name())
		}
	}
	return folders, files
}

func getCategories(db *sql.DB, rootpath string) ([]Category, map[string]int) {
	var cats []Category
	folders, _ := getContent(rootpath)
	catMap := map[string]int{}

	for i, folder := range folders {
		name := folder.Name()
		path := fmt.Sprintf("%s/%s", rootpath, folder.Name())
		catMap[name] = i

		category := Category{Name: name, Path: path}
		category.getArtist(db)
		cats = append(cats, category)
	}
	return cats, catMap
}

func main() {
	connStr := "dbname=test user=postgres host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	// categories, catMap := getCategories(db, "/media/admn/data/mp3")
	// categories[catMap["Rock"]].listContent()
	// cat := categories[catMap["Ambient"]]
	var artist string
	var name string
	var result []string

	// fmt.Println(name)
	// fmt.Println(cat.Artists[0].Albums[0].Tracks)

	// insStr := fmt.Sprintf("INSERT INTO tracks VALUES ('%[1]s', '%[2]s', '%[3]s', '%[4]s', '%[5]s') ON CONFLICT (path) DO UPDATE SET name = '%[1]s', category = '%[3]s', artist = '%[4]s', album = '%[5]s';", name, path, cat.Name, "Burzum", "album03")
	// _, err = db.Exec(insStr)
	// if err != nil {
	// 	panic(err)
	// }

	resSrt := fmt.Sprintf("SELECT DISTINCT artist, name FROM tracks WHERE category = $$%[1]s$$ AND album = $$%[1]s$$;", "Rock")
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
