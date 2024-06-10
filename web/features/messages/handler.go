package messages

import (
	"encoding/json"
	"errors"
	"messenger/data/entities"
	messagestore "messenger/data/store/message-store"
	"messenger/web/api"
	"net/http"
	"time"
)

func Handle(w http.ResponseWriter, r *http.Request, user entities.UserRef) {
	api.WithTyped(w, r, func(r *http.Request) api.ResponseAny {
		switch r.FormValue("method") {
		case "send":
			return handleSend(w, r, user)
		case "read":
			return handleRead(w, r, user)
		default:
			return api.ResponseErrPathNotFound
		}
	})
}

func handleSend(w http.ResponseWriter, r *http.Request, user entities.UserRef) api.ResponseAny {
	var form struct {
		Chat entities.ChatRef `json:"chat"`
		Text string           `json:"text"`
	}
	if e := json.Unmarshal([]byte(r.FormValue("data")), &form); e != nil {
		return api.ResponseBadRequest(e)
	}

	// Validate the message text:
	if len(form.Text) == 0 {
		return api.ResponseBadRequest(errors.New("empty message"))
	}

	//// Find recipient info from the chat info:
	//chatInfo, e := chatstore.Instance.Read(form.Chat)
	//if e != nil {
	//	return api.ResponseBadRequest(e)
	//}

	// Open requested chat:
	chat, e := messagestore.Instance.OpenChat(form.Chat, false)
	if e != nil {
		return api.ResponseInternalServerError(e)
	}
	// Send the message:
	_, e = chat.CreateMessage(func(message *entities.Message) {
		message.Sender = user
		//message.Receiver = chatInfo.OtherMember(user)
		message.Content.Text = form.Text
	})

	return api.ResponseAny{Status: 204}
}

type message struct {
	ID      entities.MessageRef `json:"id"`
	Created time.Time           `json:"created"`
	Sender  entities.UserRef    `json:"sender"`
	//Receiver entities.UserRef    `json:"receiver"`
	Content messageContent `json:"content"`
}

type messageContent struct {
	Text string `json:"text"`
}

func handleRead(w http.ResponseWriter, r *http.Request, user entities.UserRef) api.ResponseAny {
	var form struct {
		Chat entities.ChatRef `json:"chat"`

		// First and Last define the boundaries of the messages to return:
		First entities.MessageRef `json:"first"`
		Last  entities.MessageRef `json:"last"`
	}
	if e := json.Unmarshal([]byte(r.FormValue("data")), &form); e != nil {
		return api.ResponseBadRequest(e)
	}

	// Open requested chat:
	chat, e := messagestore.Instance.OpenChat(form.Chat, false)
	if e != nil {
		return api.ResponseInternalServerError(e)
	}
	// Mark all messages as read:
	e = chat.MarkAsRead(user)
	if e != nil {
		return api.ResponseInternalServerError(e)
	}
	// Get all messages:
	messages, e := chat.ReadMessages(form.First, form.Last)
	if e != nil {
		return api.ResponseInternalServerError(e)
	}

	// Convert to the web data format:
	var res []message
	for _, v := range messages {
		res = append(res, message{
			ID:      v.ID,
			Created: v.Created,
			Sender:  v.Sender,
			//Receiver: v.Receiver,
			Content: messageContent{
				Text: v.Content.Text,
			},
		})
	}

	return api.ResponseAny{Status: 200, Data: res}
}
