package authentication

import "time"

var Instance *Service

func init() {
	kp, e := diskKeysProvider()
	if e != nil {
		panic(e)
	}
	Instance = New(Options{
		GetKeys:          kp,
		CookieExpiration: time.Hour,
	})
}
