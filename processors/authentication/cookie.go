package authentication

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"messenger/data/entities"
	"net/http"
	"time"
)

func (s *Service) generateCookie(user entities.UserRef) (*http.Cookie, error) {
	var now = time.Now()
	var expires = now.Add(max(s.opts.CookieExpiration, time.Minute))

	var token = jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.RegisteredClaims{
		Subject:   string(user),
		IssuedAt:  &jwt.NumericDate{Time: now},
		ExpiresAt: &jwt.NumericDate{Time: expires},
	})
	pk, _, e := s.opts.GetKeys()
	if e != nil {
		return nil, e
	}
	tokenString, e := token.SignedString(pk)
	if e != nil {
		return nil, e
	}
	var cookie = &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expires,
	}
	return cookie, e
}

func (s *Service) validateCookieFromRequest(r *http.Request) (entities.UserRef, error) {
	cookie, e := r.Cookie("token")
	if e != nil {
		return "", e
	}

	var claims = &jwt.RegisteredClaims{}
	token, e := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (_ interface{}, e error) {
		_, pub, e := s.opts.GetKeys()
		if e != nil {
			return nil, e
		}
		return pub, nil
	})
	if e != nil {
		return "", e
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	return entities.UserRef(claims.Subject), e
}
