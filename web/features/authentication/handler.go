package authentication

import (
	"encoding/json"
	"errors"
	"messenger/data/entities"
	userstore "messenger/data/store/user-store"
	"messenger/processors/authentication"
	"messenger/web/api"
	"net/http"
	"time"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	api.WithTyped(w, r, func(r *http.Request) api.ResponseAny {
		switch r.FormValue("method") {
		case "check":
			return handleCheck(w, r)
		case "sign-in":
			return handleSignIn(w, r)
		case "sign-up":
			return handleSignUp(w, r)
		case "sign-out":
			return handleSignOut(w, r)
		default:
			return api.ResponseErrPathNotFound
		}
	})
}

// handleCheck checks whether the user is authenticated
func handleCheck(w http.ResponseWriter, r *http.Request) api.ResponseAny {
	_, e := authentication.Instance.AuthenticateByCookie(r)
	if e != nil {
		return api.ResponseBad(401, e, "Unauthorized")
	}
	return api.ResponseAny{Status: 204}
}

func handleSignIn(w http.ResponseWriter, r *http.Request) api.ResponseAny {
	var form struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if e := json.Unmarshal([]byte(r.FormValue("data")), &form); e != nil {
		return api.ResponseBadRequest(e)
	}
	// Check the provided credentials
	var user, e = userstore.Instance.Read(entities.UserRef(form.Username))
	if e != nil {
		if errors.Is(e, userstore.ErrUserNotFound) {
			return api.ResponseBad(401, e, "invalid credentials")
		} else {
			return api.ResponseBadRequest(e)
		}
	}
	if !samePassword(user.AuthMethods.PasswordHash, form.Password) {
		return api.ResponseBad(401, errors.New("invalid credentials"), "invalid credentials")
	}
	// Generate a new session
	e = authentication.Instance.SetAuthenticationCookie(w, user.ID)
	if e != nil {
		return api.ResponseInternalServerError(e)
	}
	return api.ResponseAny{Status: 204}
}

func handleSignUp(w http.ResponseWriter, r *http.Request) api.ResponseAny {
	var form struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if e := json.Unmarshal([]byte(r.FormValue("data")), &form); e != nil {
		return api.ResponseBadRequest(e)
	}

	// Validate credentials

	var validators = map[string]validationWrapper{
		"email":    {validateEmail, form.Email},
		"username": {validateUsername, form.Username},
		"password": {validatePassword, form.Password},
		"name":     {validateName, form.Name},
	}
	for key, v := range validators {
		if !v.validator(v.value) {
			return api.ResponseBad(400, errors.New(key+" is invalid"), key+" is invalid")
		}
	}

	_, e := userstore.Instance.Create(func(user *entities.User) {
		user.ID = entities.UserRef(form.Username)
		user.Credentials.Name = form.Name
		user.Credentials.Email = form.Email
		user.AuthMethods.PasswordHash = computeHash(form.Password)
	})
	if e != nil {
		if errors.Is(e, userstore.ErrUserAlreadyExists) {
			return api.ResponseBad(400, e, "user already exists")
		}

		return api.ResponseBadRequest(e)
	}
	// Write the session cookie
	e = authentication.Instance.SetAuthenticationCookie(w, entities.UserRef(form.Username))
	if e != nil {
		return api.ResponseInternalServerError(e)
	}

	return api.ResponseAny{Status: 204}
}

func handleSignOut(cookieWriter http.ResponseWriter, r *http.Request) api.ResponseAny {
	http.SetCookie(cookieWriter, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})
	return api.ResponseAny{Status: 204}
}
