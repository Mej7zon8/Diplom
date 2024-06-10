package chats

import (
	"fmt"
	"messenger/data/entities"
	chatstore "messenger/data/store/chat-store"
	messagestore "messenger/data/store/message-store"
	userstore "messenger/data/store/user-store"
	"slices"
	"time"
)

type chat struct {
	ID           entities.ChatRef `json:"id"`
	Name         string           `json:"name"`
	Participants []Participant    `json:"participants"`
	LastMessage  message          `json:"last_message"`
	IsGroup      bool             `json:"is_group"`
}

type Participant struct {
	ID   entities.UserRef `json:"id"`
	Name string           `json:"name"`
}

type message struct {
	ID       entities.MessageRef `json:"id"`
	Created  time.Time           `json:"created"`
	Sender   entities.UserRef    `json:"sender"`
	Receiver entities.UserRef    `json:"receiver"`
	Content  messageContent      `json:"content"`
}

type messageContent struct {
	Text string `json:"text"`
}

func buildUserChats(user entities.UserRef) (res []*chat, err error) {
	// List all chats for the user:
	data, err := chatstore.Instance.GetUserChats(user)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		// Get other participants' info:
		var participants []Participant
		for _, p := range v.Participants {
			if p != user {
				otherUser, e := userstore.Instance.Read(p)
				if e != nil {
					return nil, e
				}
				participants = append(participants, Participant{
					ID:   otherUser.ID,
					Name: otherUser.Credentials.Name,
				})
			}
		}
		var chatName string
		if !v.IsGroup {
			chatName = participants[0].Name
		} else {
			if v.Name != "" {
				chatName = v.Name[:min(len(v.Name), 20)]
			} else {
				chatName = fmt.Sprintf("%v участников", len(participants)+1)
			}
		}

		// Build the chat statistics:
		stats, e := getChatStats(v.ID, user)
		if e != nil {
			return nil, e
		}

		res = append(res, &chat{
			ID:           v.ID,
			Name:         chatName,
			Participants: participants,
			IsGroup:      v.IsGroup,
			LastMessage: message{
				ID:      stats.LastMessage.ID,
				Created: stats.LastMessage.Created,
				Sender:  stats.LastMessage.Sender,
				//Receiver: stats.LastMessage.Receiver,
				Content: messageContent{
					Text: stats.LastMessage.Content.Text,
				},
			},
		})
	}

	// Sort the chats by the last message date:
	slices.SortFunc(res, func(a, b *chat) int {
		if b.LastMessage.Created.After(a.LastMessage.Created) {
			return 1
		}
		if b.LastMessage.Created.Before(a.LastMessage.Created) {
			return -1
		}
		return 0
	})
	return res, nil
}

func getChatStats(chat entities.ChatRef, user entities.UserRef) (entities.ChatStats, error) {
	// Open the chat in the message store:
	c, e := messagestore.Instance.OpenChat(chat, false)
	if e != nil {
		return entities.ChatStats{}, e
	}

	// Build the statistics:
	return c.BuildStatistic(user)
}
