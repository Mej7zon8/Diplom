package chat_store

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

const tableName = "chats"

func New() *Store {
	var table = database.ChatDatabase.Table(tableName)
	if err := table.AutoMigrate(&entities.ChatInfo{}); err != nil {
		panic(err)
	}
	return &Store{
		table: table,
	}
}
