package authentication

import (
	"github.com/k773/utils"
	"net/mail"
)

type validationWrapper struct {
	validator func(string) bool
	value     string
}

var allowedUsernameRunes = utils.Slice2HasMap([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))

func validateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	for _, r := range username {
		if _, h := allowedUsernameRunes[r]; !h {
			return false
		}
	}
	return true
}

func validateEmail(data string) bool {
	_, e := mail.ParseAddress(data)
	return e == nil
}

func validatePassword(data string) bool {
	return len(data) >= 8
}

func validateName(data string) bool {
	return len(data) >= 3 && len(data) <= 50
}
