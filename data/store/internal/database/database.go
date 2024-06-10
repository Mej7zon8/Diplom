package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type Database = gorm.DB

var (
	UserDatabase *Database
	ChatDatabase *Database
)

func init() {
	var err error

	UserDatabase, err = open("users.sqlite3")
	if err != nil {
		panic(err)
	}
	ChatDatabase, err = open("chats.sqlite3")
	if err != nil {
		panic(err)
	}
}

func open(path string) (*Database, error) {
	_ = os.MkdirAll("program-data/db/", 0600)

	return gorm.Open(sqlite.Open(filepath.Join("program-data/db/", path)), &gorm.Config{})
}
