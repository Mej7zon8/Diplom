package message_store

import (
	"errors"
	"messenger/data/entities"
	"messenger/data/store/internal/database"
)

var Instance *Store

func init() {
	Instance = New()
}

type Store struct {
	database *database.Database
}

type Chat struct {
	table *database.Database
}

func New() *Store {
	return &Store{
		database: database.ChatDatabase,
	}
}

func (s *Store) OpenChat(chatRef entities.ChatRef, allowCreate bool) (*Chat, error) {
	var chat *Chat
	var tableName = "chat_" + string(chatRef)

	//e := s.database.Transaction(func(tx *gorm.DB) error {
	if !allowCreate {
		// Check if the table does not exist.
		if !s.database.Migrator().HasTable(tableName) {
			return nil, errors.New("chat does not exist")
		}
	}
	chat = &Chat{
		table: s.database.Table(tableName),
	}
	return chat, chat.table.AutoMigrate(&entities.Message{})
	//})
	//return chat, e
}
