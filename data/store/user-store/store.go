package user_store

import (
	"messenger/data/entities"
	"messenger/data/store/internal/database"
)

var Instance *Store

func init() {
	Instance = New()
}

type Store struct {
	table *database.Database
}

func New() *Store {
	var table = database.UserDatabase.Table("users")
	if err := table.AutoMigrate(&entities.User{}); err != nil {
		panic(err)
	}
	return &Store{
		table: table,
	}
}
