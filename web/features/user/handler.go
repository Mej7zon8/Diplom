package user

import (
	"encoding/json"
	"errors"
	"messenger/data/entities"
	messagestore "messenger/data/store/message-store"
	userstore "messenger/data/store/user-store"
	"messenger/web/api"
	"net/http"
	"time"
)

func Handle(w http.ResponseWriter, r *http.Request, user entities.UserRef) {
	api.WithTyped(w, r, func(r *http.Request) api.ResponseAny {
		switch r.FormValue("method") {
		case "exist":
			return handleExist(w, r)
		default:
			return api.ResponseErrPathNotFound
		}
	})
}

func handleExist(w http.ResponseWriter, r *http.Request) api.ResponseAny {
	var form struct {
		User entities.UserRef `json:"user"`
	}
	if e := json.Unmarshal([]byte(r.FormValue("data")), &form); e != nil {
		return api.ResponseBadRequest(e)
	}

	_, e := userstore.Instance.Read(form.User)
	if e != nil && !errors.Is(e, userstore.ErrUserNotFound) {
		return api.ResponseInternalServerError(e)
	}
	return api.ResponseAny{Status: 200, Data: e == nil}
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
