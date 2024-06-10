package categorize

import (
	"context"
	"errors"
	"messenger/data/entities"
	messagestore "messenger/data/store/message-store"
	userstore "messenger/data/store/user-store"
	messagecategorization "messenger/processors/message-categorization"
)

// chatWorker is a worker that processes all uncategorized messages in the chat.
type chatWorker struct {
	chatRef entities.ChatRef
	chat    *messagestore.Chat
}

// newChatWorker creates a new chat worker.
func newChatWorker(chatRef entities.ChatRef) *chatWorker {
	return &chatWorker{chatRef: chatRef}
}

// Run processes all uncategorized messages in the chat.
func (c *chatWorker) Run() (e error) {
	// open the chat:
	c.chat, e = messagestore.Instance.OpenChat(c.chatRef, false)
	if e != nil {
		return e
	}

	// process the messages until there are no more messages to process:
	var next bool
	for {
		next, e = c.next()
		if e != nil || !next {
			return e
		}
	}
}

// next processes the next batch of uncategorized messages in the chat.
// returns:
//   - next: true if there are more messages to process, false otherwise.
//     should be ignored if an error occurs.
//   - e: the error that occurred during the processing of the messages.
func (c *chatWorker) next() (next bool, e error) {
	// Get the next batch of uncategorized messages.
	messages, e := c.chat.GetNextUncategorized(10)
	if e != nil {
		return false, e
	}
	if len(messages) == 0 {
		return false, nil
	}

	// Process each message:
	var accum = newAccumulator()
	for _, message := range messages {
		if len(message.Content.Text) < 3 { // Ignore too short messages.
			continue
		}

		var res messagecategorization.ProcessRes
		res, e = messagecategorization.Instance.Process(context.Background(), message.Content.Text)
		if e != nil { // we still want to record the collected stats
			break
		}

		// Write the result to the message.
		e = c.chat.UpdateMessage(message.ID, func(m *entities.Message) {
			m.Meta.Categories.Processed = true
			m.Meta.Categories.Values = res.Result
		})
		if e != nil { // we still want to record the collected stats
			break
		}
		for label, value := range res.Result {
			accum.add(message.Sender, label, value)
		}
	}

	// Save the collected stats to the user's profile:
	for userRef, label2value := range accum.user2label2value {
		// Update the user's profile.
		err := userstore.Instance.Update(userRef, func(u *entities.User) {
			for label, value := range label2value {
				u.Meta.Categories.Values[label] += value
			}
		})
		// Do not interrupt the update process if an error occurs.
		e = errors.Join(e, err)
	}
	return len(accum.user2label2value) != 0, e
}
