package authentication

import (
	"messenger/data/entities"
	"net/http"
)

func (s *Service) AuthenticateByCookie(r *http.Request) (entities.UserRef, error) {
	return s.validateCookieFromRequest(r)
}

func (s *Service) SetAuthenticationCookie(w http.ResponseWriter, user entities.UserRef) error {
	cookie, e := s.generateCookie(user)
	if e != nil {
		return e
	}
	http.SetCookie(w, cookie)
	return nil
}
