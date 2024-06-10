package authentication

import (
	"crypto/sha512"
	"encoding/hex"
)

func samePassword(passwordHash, password string) bool {
	return computeHash(password) == passwordHash
}

func computeHash(password string) string {
	var a = sha512.Sum512([]byte(password))
	return hex.EncodeToString(a[:])
}
