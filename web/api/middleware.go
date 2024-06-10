package api

import (
	"encoding/json"
	"log/slog"
	"messenger/data/entities"
	"messenger/processors/authentication"
	"net/http"
)

func AuthMiddleware(w http.ResponseWriter, r *http.Request, next func(w http.ResponseWriter, r *http.Request, user entities.UserRef)) {
	user, e := authentication.Instance.AuthenticateByCookie(r)
	if e != nil {
		_ = json.NewEncoder(w).Encode(Response[any]{Status: 401, Error: "Unauthorized"})
		return
	}
	next(w, r, user)
}

func WithTyped[T any](w http.ResponseWriter, r *http.Request, next func(r *http.Request) Response[T]) {
	data := next(r)
	if data.Error != "" {
		var l = slog.With("module", "api", "identifiers", extractIdentifiers(r), "ux-error", data.Error, "status", data.Status)
		if data.RealError != nil {
			l = l.With("real error", data.RealError)
		}
		l.Debug("failed")
	}
	_ = json.NewEncoder(w).Encode(data)
}

func WithTypedWithUser[T any](w http.ResponseWriter, r *http.Request, next func(r *http.Request, user entities.UserRef) Response[T]) {
	AuthMiddleware(w, r, func(w http.ResponseWriter, r *http.Request, user entities.UserRef) {
		data := next(r, user)
		_ = json.NewEncoder(w).Encode(data)
	})
}

func extractIdentifiers(r *http.Request) []string {
	var reqIdentifiers []string
	for _, possible := range []string{"module", "submodule", "method", "uuid", "groupId", "userId", "user"} {
		if v := r.FormValue(possible); v != "" {
			// Clip the value to 255 symbols
			if len(v) > 255 {
				v = v[:255] + "..."
			}
			reqIdentifiers = append(reqIdentifiers, possible+": "+v)
		}
	}
	return reqIdentifiers
}
