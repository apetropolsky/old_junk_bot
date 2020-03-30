package library

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
)

// Track structure
type Track struct {
	Name, Path, Category, Album, Artist string
}

// GetContent of library
func GetContent(db *sql.DB, rootpath string) ([]Track, error) {
	var tracks []Track

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			category := strings.Split(path, "/")[len(strings.Split(rootpath, "/"))]
			artist := strings.Split(path, "/")[len(strings.Split(rootpath, "/"))+1]
			album := filepath.Base(filepath.Dir(path))
			name := filepath.Base(path)

			track := Track{
				Name:     name,
				Path:     path,
				Category: category,
				Artist:   artist,
				Album:    album,
			}

			tracks = append(tracks, track)
		}
		return nil
	})
	return tracks, err
}
