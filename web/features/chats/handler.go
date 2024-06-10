package chats

import (
	"encoding/json"
	"errors"
	"messenger/data/entities"
	chatstore "messenger/data/store/chat-store"
	"messenger/data/store/message-store"
	userstore "messenger/data/store/user-store"
	"messenger/web/api"
	"net/http"
	"slices"
	"strconv"
)

func Handle(w http.ResponseWriter, r *http.Request, user entities.UserRef) {
	api.WithTyped(w, r, func(r *http.Request) api.ResponseAny {
		switch r.FormValue("method") {
		case "create":
			return handleCreate(w, r, user)
		case "list":
			return handleList(w, r, user)
		case "leave":
			return handleLeave(w, r, user)
		default:
			return api.ResponseErrPathNotFound
		}
	})
}

func handleCreate(w http.ResponseWriter, r *http.Request, user entities.UserRef) api.ResponseAny {
	var request struct {
		Name         string             `json:"name"`
		Participants []entities.UserRef `json:"participants"`
	}
	if e := json.Unmarshal([]byte(r.FormValue("data")), &request); e != nil {
		return api.ResponseBadRequest(e)
	}

	if len(request.Participants) == 0 || len(request.Participants) > 10 {
		return api.ResponseBadRequest(errors.New("invalid number of participants: " + strconv.Itoa(len(request.Participants))))
	}
	for _, participant := range request.Participants {
		// Disallow creating chats with the same user:
		if participant == user {
			return api.ResponseBad(400, errors.New("cannot create chat with self"), "cannot create chat with self")
		}
		// Check that the user exists:
		_, e := userstore.Instance.Read(participant)
		if e != nil {
			return api.ResponseBad(400, e, "user not found")
		}
	}
	// Add user to the list of participants:
	request.Participants = append(request.Participants, user)

	// Create the chat in the chat store
	chat, e := chatstore.Instance.Create(func(chat *entities.ChatInfo) {
		chat.Name = request.Name
		chat.IsGroup = len(request.Participants) > 1
		chat.Participants = request.Participants
	})
	if e != nil {
		if errors.Is(e, chatstore.ErrChatExists) {
			return api.ResponseBad(400, e, "chat already exists")
		}
		return api.ResponseInternalServerError(e)
	}
	// Create the chat in the message store
	_, e = message_store.Instance.OpenChat(chat.ID, true)
	return api.ResponseAny{Status: 204}
}

func handleList(w http.ResponseWriter, r *http.Request, user entities.UserRef) api.ResponseAny {
	res, e := buildUserChats(user)
	if e != nil {
		return api.ResponseInternalServerError(e)
	}
	return api.ResponseAny{Status: 200, Data: res}
}

func handleLeave(w http.ResponseWriter, r *http.Request, user entities.UserRef) api.ResponseAny {
	var request struct {
		Chat entities.ChatRef `json:"chat"`
	}

	if e := json.Unmarshal([]byte(r.FormValue("data")), &request); e != nil {
		return api.ResponseBadRequest(e)
	}

	chatInfo, e := chatstore.Instance.Read(request.Chat)
	if e != nil {
		return api.ResponseInternalServerError(e)
	}
	if !chatInfo.IsGroup {
		return api.ResponseBad(400, errors.New("cannot leave a private chat"), "cannot leave a private chat")
	}

	// Remove the user from the list of participants
	e = chatstore.Instance.Update(request.Chat, func(chat *entities.ChatInfo) {
		chat.Participants = slices.DeleteFunc(chat.Participants, func(ref entities.UserRef) bool {
			return ref == user
		})
	})
	if e != nil {
		return api.ResponseInternalServerError(e)
	}
	return api.ResponseAny{Status: 204}
}
