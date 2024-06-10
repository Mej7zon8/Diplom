package authentication

import (
	"time"
)

type Service struct {
	opts Options
}

type Options struct {
	GetKeys          keysProvider
	CookieExpiration time.Duration
}

func New(opts Options) *Service {
	return &Service{opts: opts}
}
