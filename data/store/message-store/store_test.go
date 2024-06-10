package message_store

import (
	"fmt"
	"messenger/data/entities"
	"testing"
)

func TestStore(t *testing.T) {
	var store = New()

	// Try and open a chat that does not exist without allowing its creation.
	{
		_, e := store.OpenChat("test-nonexistent", false)
		if e == nil {
			t.Fatal("expected an error")
		}
	}
	// Try and open a chat that does not exist and allow its creation.
	{
		_, e := store.OpenChat("test", true)
		if e != nil {
			t.Fatal(e)
		}
	}

	// Try and open a chat that exists.
	chat, e := store.OpenChat("test", false)
	if e != nil {
		t.Fatal(e)
	}

	// Try and read messages from the chat. The result should be empty.
	{
		_, e = chat.ReadMessages(-1, -1)
		if e != nil {
			t.Fatal(e)
		}
	}
	// Try and insert 3 messages into the chat.
	{
		for i := 0; i < 3; i++ {
			_, e = chat.CreateMessage(func(m *entities.Message) {
				m.Content.Text = fmt.Sprintf("message %d", i+1)
			})
			if e != nil {
				t.Fatal(e)
			}
		}
	}
	// Try and read messages from the chat. The result should contain 3 messages ordered from the latest to the oldest.
	{
		var messages, e = chat.ReadMessages(-1, -1)
		if e != nil {
			t.Fatal(e)
		}
		if len(messages) != 3 {
			t.Fatal("expected 3 messages")
		}
		for i, m := range messages {
			if m.Content.Text != fmt.Sprintf("message %d", 3-i) {
				t.Fatalf("expected message %v, got %v", fmt.Sprintf("message %d", 2-i), m.Content.Text)
			}
		}
	}
	// Try and read a message with id 1 from the chat.
	{
		var messages, e = chat.ReadMessages(1, 1)
		if e != nil {
			t.Fatal(e)
		}
		if len(messages) != 1 {
			t.Fatal("expected 1 message")
		}
		if messages[0].Content.Text != "message 1" {
			t.Fatal("expected message 1, got", messages[0].Content.Text)
		}
	}
	// Update the message with id 1.
	{
		e = chat.UpdateMessage(1, func(m *entities.Message) {
			m.Content.Text = "updated message"
		})
		if e != nil {
			t.Fatal(e)
		}
	}
	// Try and read the updated message with id 1 from the chat.
	{
		var messages, e = chat.ReadMessages(1, 1)
		if e != nil {
			t.Fatal(e)
		}
		if len(messages) != 1 {
			t.Fatal("expected 1 message")
		}
		if messages[0].Content.Text != "updated message" {
			t.Fatal("expected updated message")
		}
	}
	// Delete the message with id 1.
	{
		e = chat.DeleteMessage(1)
		if e != nil {
			t.Fatal(e)
		}
	}
	// Try and read the deleted message with id 1 from the chat.
	{
		var messages, e = chat.ReadMessages(1, 1)
		if e != nil {
			t.Fatal(e)
		}
		if len(messages) != 0 {
			t.Fatal("expected 0 messages")
		}
	}
}
