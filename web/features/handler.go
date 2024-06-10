package features

import (
	"encoding/json"
	"messenger/data/entities"
	"messenger/web/api"
	"messenger/web/features/advertisement"
	"messenger/web/features/authentication"
	"messenger/web/features/chats"
	"messenger/web/features/messages"
	userfeature "messenger/web/features/user"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	var module = r.FormValue("module")
	// We need to skip the authentication for some modules:
	switch module {
	case "authentication":
		authentication.Handle(w, r)
		return
	}
	// Proceed with the authentication:
	api.AuthMiddleware(w, r, func(w http.ResponseWriter, r *http.Request, user entities.UserRef) {
		// Handle the authenticated request:
		switch module {
		case "chats":
			chats.Handle(w, r, user)
		case "messages":
			messages.Handle(w, r, user)
		case "user":
			userfeature.Handle(w, r, user)
		case "advertisement":
			advertisement.Handle(w, r, user)
		default:
			_ = json.NewEncoder(w).Encode(api.ResponseErrPathNotFound)
		}
	})
}
