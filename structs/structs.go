package structs

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
)

// Category on the top level of library
type Category struct {
	Name, Path string
	Artists    []Artist
}

// GetArtist from Category
func (c *Category) GetArtist(db *sql.DB) {
	folders, _ := GetContent(c.Path)

	for _, folder := range folders {
		name := folder.Name()
		path := fmt.Sprintf("%s/%s", c.Path, name)

		artist := Artist{
			Name:     name,
			Path:     path,
			Category: c.Name,
		}

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

func (a *Artist) getAlbums(db *sql.DB) {
	a.Tracks = getTracks(
		db,
		a.Path,
		a.Category,
		a.Category, // For unsorted files use a category name as an album name
		a.Name,
	)

	albums, _ := GetContent(a.Path)

	for _, alb := range albums {
		name := alb.Name()
		path := fmt.Sprintf("%s/%s", a.Path, alb.Name())

		album := Album{
			Name:     name,
			Path:     path,
			Category: a.Category,
			Artist:   a.Name,
		}

		album.Tracks = getTracks(
			db,
			album.Path,
			a.Category,
			album.Name,
			a.Name,
		)

		a.Albums = append(a.Albums, album)
	}
}

// Album in artist folder
type Album struct {
	Name, Path, Category, Artist string
	Tracks                       []Track
}

// Track structure
type Track struct {
	Name, Path, Category, Album, Artist string
}

func getTracks(db *sql.DB,
	rootpath string,
	category string,
	album string,
	artist string) []Track {

	var tracks []Track
	_, files := GetContent(rootpath)

	for _, file := range files {
		path := fmt.Sprintf("%s/%s", rootpath, file)

		track := Track{
			Name:     file,
			Path:     path,
			Category: category,
			Album:    album,
			Artist:   artist,
		}

		tracks = append(tracks, track)
		insertTrack(db, track)
	}
	return tracks
}

func insertTrack(db *sql.DB, track Track) {
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

// GetContent of directory
func GetContent(rootpath string) ([]os.FileInfo, []string) {
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
