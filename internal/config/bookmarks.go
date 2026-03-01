package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type BookmarksConfig struct {
	Bookmarks BookmarksType `toml:"bookmarks"`
}

type BookmarksType map[string]string

var bookmarks = BookmarksType{}

func LoadBookmarks() error {
	bmDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	bmPath := filepath.Join(bmDir, "duvet", "bookmarks.toml")
	if _, err := os.Stat(bmPath); err != nil && os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(bmPath)
	if err != nil {
		return err
	}

	var bmc BookmarksConfig
	if err := toml.Unmarshal(data, &bmc); err != nil {
		return err
	}

	if bmc.Bookmarks != nil {
		bookmarks = bmc.Bookmarks
	}

	return nil
}

func SaveBookmarks() error {
	bmDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	dirPath := filepath.Join(bmDir, "duvet")
	bmPath := filepath.Join(dirPath, "bookmarks.toml")

	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return err
	}

	f, err := os.Create(bmPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(BookmarksConfig{Bookmarks: bookmarks})
}

func GetBookmark(name string) (string, bool) {
	path, ok := bookmarks[name]
	return path, ok
}

func GetBookmarks() BookmarksType {
	return bookmarks
}

func SetBookmark(name, path string) error {
	bookmarks[name] = path
	return SaveBookmarks()
}

func DeleteBookmark(name string) error {
	delete(bookmarks, name)
	return SaveBookmarks()
}
