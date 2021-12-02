package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

type Cipher interface {
	Encrypt(message []byte, privateKey *rsa.PrivateKey, receiverPublicKey *rsa.PublicKey) ([]byte, error)
	Decrypt(message []byte, privateKey *rsa.PrivateKey, senderPublicKey *rsa.PublicKey) ([]byte, error)
}

func NewCipher() Cipher {
	return &rsaOaepPss{
		hashType: crypto.SHA256,
		hashFunc: func(in []byte) []byte {
			bytes := sha256.Sum256(in)
			return bytes[:]
		},
		rng: rand.Reader,
	}
}
