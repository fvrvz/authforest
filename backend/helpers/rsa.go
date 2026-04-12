package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"sync"

	"github.com/fvrvz/gologger"
)

var (
	rsaPrivateKey *rsa.PrivateKey
	rsaOnce       sync.Once
	rsaKeyID      = "authforest-key-1"
)

func InitRSAKey(keyPath string) {
	rsaOnce.Do(func() {
		if _, err := os.Stat(keyPath); err == nil {
			data, err := os.ReadFile(keyPath)
			if err != nil {
				log.Fatalf("Failed to read RSA key file: %v", err)
			}
			block, _ := pem.Decode(data)
			if block == nil {
				log.Fatal("Failed to decode PEM block from RSA key file")
			}
			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				log.Fatalf("Failed to parse RSA private key: %v", err)
			}
			rsaPrivateKey = key
			gologger.INFO("RSA private key loaded from %s", keyPath)
		} else {
			gologger.INFO("Generating new RSA key pair (2048-bit)")
			key, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				log.Fatalf("Failed to generate RSA key: %v", err)
			}
			rsaPrivateKey = key

			pemData := pem.EncodeToMemory(&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(key),
			})
			if err := os.WriteFile(keyPath, pemData, 0600); err != nil {
				gologger.WARN("Failed to persist RSA key to %s: %v (key will be regenerated on restart)", keyPath, err)
			} else {
				gologger.INFO("RSA private key persisted to %s", keyPath)
			}
		}
	})
}

func GetRSAPrivateKey() *rsa.PrivateKey {
	return rsaPrivateKey
}

func GetRSAPublicKey() *rsa.PublicKey {
	return &rsaPrivateKey.PublicKey
}

func GetRSAKeyID() string {
	return rsaKeyID
}
