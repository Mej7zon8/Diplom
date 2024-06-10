package message_store

import (
	"messenger/data/entities"
	"messenger/data/store/internal/database"
)

// ReadMessages reads messages from the chat in the range [first, last] ordered from the oldest to the newest.
// If a range element equals to -1, the range is unbounded in that direction.
func (c *Chat) ReadMessages(first, last entities.MessageRef) (messages []entities.Message, err error) {
	err = c.table.Transaction(func(tx *database.Database) error {
		var query = tx.Order("id DESC")
		if first != -1 {
			query = query.Where("id >= ?", first)
		}
		if last != -1 {
			query = query.Where("id <= ?", last)
		}
		return query.Find(&messages).Error
	})
	return
}

// GetNextUncategorized retrieves the next batch of uncategorized messages from the chat.
func (c *Chat) GetNextUncategorized(batchSize int) (messages []entities.Message, err error) {
	err = c.table.Transaction(func(tx *database.Database) error {
		return tx.Where("meta_categories_processed = ?", false).Limit(batchSize).Find(&messages).Error
	})
	return
}

func (c *Chat) BuildStatistic(requester entities.UserRef) (res entities.ChatStats, err error) {
	err = c.table.Transaction(func(tx *database.Database) error {
		// Count the number of unread messages in the chat for the user.
		var count int64
		e := tx.Where("sender != ?", requester).Count(&count).Error
		if e != nil {
			return e
		}
		//res.UnreadCount = int(count)
		// Retrieve the last message in the chat.
		e = tx.Order("id DESC").First(&res.LastMessage).Error
		// Ignore the record not found error.
		if database.IsErrorRecordNotFound(e) {
			return nil
		}
		return e
	})
	return res, err
}

func (c *Chat) MarkAsRead(user entities.UserRef) error {
	// no-op
	return nil
	//return c.table.Transaction(func(tx *database.Database) error {
	//	return tx.Model(&entities.Message{}).Where("receiver = ? AND read = ?", user, false).Update("read", true).Error
	//})
}

// CreateMessage creates a new message in the chat.
func (c *Chat) CreateMessage(update func(*entities.Message)) (entities.Message, error) {
	var message = entities.NewMessage()
	update(message)
	var e = c.table.Transaction(func(tx *database.Database) error { return tx.Create(message).Error })
	return *message, e
}

// UpdateMessage updates the message in the chat.
func (c *Chat) UpdateMessage(messageRef entities.MessageRef, update func(*entities.Message)) error {
	return c.table.Transaction(func(tx *database.Database) error {
		var message entities.Message
		err := tx.First(&message, messageRef).Error
		if err != nil {
			return err
		}
		update(&message)
		return tx.Save(&message).Error
	})
}

// DeleteMessage deletes the message from the chat.
func (c *Chat) DeleteMessage(messageRef entities.MessageRef) error {
	return c.table.Transaction(func(tx *database.Database) error {
		return tx.Delete(&entities.Message{}, messageRef).Error
	})
}
