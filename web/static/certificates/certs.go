package certificates

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

// PublicKey retrieves public key from private key
func PublicKey(privateKey any) any {
	switch k := privateKey.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		panic("Unknown private key format")
	}
}

// PemBlockFromKey creates pem block from the provided private key
func PemBlockFromKey(privateKey interface{}) *pem.Block {
	switch k := privateKey.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			log.Fatalf("Unable to marshal ECDSA private key: %v", err)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		panic("Unknown private key format")
	}
}

func GenerateCertificatesIfNotFound(certPath, keyPath string) {
	if !FileExist(certPath) || !FileExist(keyPath) {
		println("certs.go: generating certificates")
		_ = os.Remove(certPath)
		_ = os.Remove(keyPath)
	} else {
		return
	}

	// privateKey, err := rsa.GenerateKey(rand.Reader, *rsaBits)
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)

	if err != nil {
		log.Fatal(err)
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{""},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 1800),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, PublicKey(privateKey), privateKey)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}
	out := &bytes.Buffer{}
	panicOnError(pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}))
	panicOnError(os.WriteFile(certPath, out.Bytes(), 0600))
	out.Reset()
	panicOnError(pem.Encode(out, PemBlockFromKey(privateKey)))
	panicOnError(os.WriteFile(keyPath, out.Bytes(), 0600))
}
