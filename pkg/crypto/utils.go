package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
)

var (
	b64Decode = base64.StdEncoding.DecodeString
	b64Encode = base64.StdEncoding.EncodeToString
)

func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return privateKey, &privateKey.PublicKey
}

func EncodeKey(key *rsa.PublicKey) (string, error) {
	keyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	return b64Encode(keyBytes), nil
}

func DecodeKey(publicKeyBase64 string) (*rsa.PublicKey, error) {
	b, err := b64Decode(publicKeyBase64)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not a key")
	}
	return rsaKey, nil
}
