package entities

import (
	"github.com/google/uuid"
	"messenger/data/common"
	"time"
)

type ChatRef string

type ChatInfo struct {
	ID      ChatRef
	Created time.Time

	Name         string
	IsGroup      bool
	Participants common.JsonArray[UserRef]
}

type ChatStats struct {
	LastMessage Message
}

// Implementation:

func NewChat() *ChatInfo {
	return &ChatInfo{
		ID:      ChatRef(uuid.NewString()),
		Created: time.Now(),
	}
}

func (c *ChatInfo) OtherMembers(exclude UserRef) []UserRef {
	var res []UserRef
	for _, v := range c.Participants {
		if v != exclude {
			res = append(res, v)
		}
	}
	return res
}
