package authentication

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"os"
	"slices"
)

type keysProvider func() (ed25519.PrivateKey, ed25519.PublicKey, error)

func diskKeysProvider() (keysProvider, error) {
	const dir = "program-data/authentication-keys"
	var pkPath = dir + "/pk_ed25519"

	_ = os.MkdirAll(dir, 0600)

	// read private key bytes from file
	keyBytes, e := os.ReadFile(pkPath)
	if e != nil && !errors.Is(e, os.ErrNotExist) {
		return nil, e
	}
	// Generate keypair if not exist
	if len(keyBytes) == 0 {
		_, keyBytes, e = ed25519.GenerateKey(rand.Reader)
		if e != nil {
			return nil, e
		}
		// Save private key bytes to file
		e = os.WriteFile(pkPath, keyBytes, 0600)
		if e != nil {
			return nil, e
		}
	}
	var pk = ed25519.PrivateKey(keyBytes)
	var pub = pk.Public().(ed25519.PublicKey)
	return func() (ed25519.PrivateKey, ed25519.PublicKey, error) {
		return slices.Clone(pk), slices.Clone(pub), nil
	}, nil
}
