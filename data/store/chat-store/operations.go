package chat_store

import (
	"encoding/json"
	"errors"
	"messenger/data/entities"
	"messenger/data/store/internal/database"
)

var ErrChatExists = errors.New("chat already exists")

// Create creates a new record in the database for the specified chat.
func (s *Store) Create(update func(*entities.ChatInfo)) (entities.ChatInfo, error) {
	var chat = entities.NewChat()
	update(chat)
	var e = s.table.Transaction(func(tx *database.Database) error {
		// Disallow creating a duplicate chat for two users.
		if len(chat.Participants) == 2 {
			// Here goes a hacky solution to check if the chat already exists.
			// Sqlite does not support array search, and doing it in Go is inefficient.
			// Therefore, we have to search for marshaled values in the database.
			var assertNotExists = func(variant []entities.UserRef) (e error) {
				variantJson, e := json.Marshal(variant)
				if e != nil {
					return e
				}
				var count int64
				e = tx.Where("participants = ?", variantJson).Count(&count).Error
				if e != nil {
					return
				}
				if count > 0 {
					return ErrChatExists
				}
				return
			}
			e := assertNotExists([]entities.UserRef{chat.Participants[0], chat.Participants[1]})
			if e != nil {
				return e
			}
			e = assertNotExists([]entities.UserRef{chat.Participants[1], chat.Participants[0]})
			if e != nil {
				return e
			}
		}
		return tx.Create(chat).Error
	})
	return *chat, e
}

// Read retrieves the record from the database for the specified chat.
func (s *Store) Read(chatRef entities.ChatRef) (entities.ChatInfo, error) {
	var chat entities.ChatInfo
	return chat, s.table.Transaction(func(tx *database.Database) error {
		return tx.First(&chat, chatRef).Error
	})
}

func (s *Store) GetAllChats() ([]entities.ChatInfo, error) {
	var chats []entities.ChatInfo
	return chats, s.table.Transaction(func(tx *database.Database) error {
		return tx.Find(&chats).Error
	})
}

// GetUserChats retrieves the list of chats where the user is a member.
// The returned list is to be assumed as unordered.
func (s *Store) GetUserChats(user entities.UserRef) ([]entities.ChatInfo, error) {
	var chats []entities.ChatInfo
	return chats, s.table.Transaction(func(tx *database.Database) error {
		// Select only the chats where the user is a member.
		// Due to the use of json_each, the queried fields must be explicitly listed.
		var query = tx.Raw("SELECT chats.id, created, participants, name, is_group FROM chats, json_each(chats.participants) WHERE json_each.value = ?", user)

		return query.Find(&chats).Error
	})
}

// Update updates the record in the database for the specified chat.
func (s *Store) Update(chatRef entities.ChatRef, update func(*entities.ChatInfo)) error {
	return s.table.Transaction(func(tx *database.Database) error {
		var chat, err = s.Read(chatRef)
		if err != nil {
			return err
		}
		update(&chat)
		return tx.Save(&chat).Error
	})
}

func (s *Store) Delete(chatRef entities.ChatRef) error {
	return s.table.Transaction(func(tx *database.Database) error {
		return tx.Delete(&entities.ChatInfo{}, chatRef).Error
	})
}
