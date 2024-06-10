package entities

import (
	"messenger/data/common"
	"time"
)

type MessageRef int

type Message struct {
	ID      MessageRef `gorm:"primaryKey,autoIncrement"`
	Created time.Time

	Sender UserRef
	//Receiver UserRef
	Content MessageContent `gorm:"embedded;embeddedPrefix:content_"`
	Meta    MessageMeta    `gorm:"embedded;embeddedPrefix:meta_"`
}

type MessageContent struct {
	Text string
}

type MessageMeta struct {
	Categories MessageCategories `gorm:"embedded;embeddedPrefix:categories_"`
}

type MessageCategories struct {
	Processed bool
	Values    common.JsonMap[string, float64]
}

func NewMessage() *Message {
	return &Message{
		ID:      0, // The id will be automatically assigned by the database.
		Created: time.Now(),
		Meta: MessageMeta{
			Categories: MessageCategories{
				Values: common.JsonMap[string, float64]{},
			},
		},
	}
}
