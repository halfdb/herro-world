package crypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/gob"
	"io"
)

type rsaOaepPss struct {
	hashType crypto.Hash
	hashFunc func([]byte) []byte
	rng      io.Reader
}

type signedMessage struct {
	CipherText []byte
	Signature  []byte
}

func marshalSignedMessage(message *signedMessage) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(message)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func unmarshalSignedMessage(b []byte) (*signedMessage, error) {
	message := &signedMessage{}
	buffer := bytes.Buffer{}
	buffer.Write(b)
	d := gob.NewDecoder(&buffer)
	err := d.Decode(message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (r *rsaOaepPss) Encrypt(message []byte, privateKey *rsa.PrivateKey, receiverPublicKey *rsa.PublicKey) ([]byte, error) {
	digest := r.hashFunc(message)
	signature, err := rsa.SignPSS(r.rng, privateKey, r.hashType, digest, nil)
	if err != nil {
		return nil, err
	}
	cipherText, err := rsa.EncryptOAEP(r.hashType.New(), rand.Reader, receiverPublicKey, message, nil)
	if err != nil {
		return nil, err
	}
	payload := &signedMessage{
		CipherText: cipherText,
		Signature:  signature,
	}

	marshaled, err := marshalSignedMessage(payload)
	if err != nil {
		return nil, err
	}

	return marshaled, nil
}

func (r *rsaOaepPss) Decrypt(message []byte, privateKey *rsa.PrivateKey, senderPublicKey *rsa.PublicKey) ([]byte, error) {
	signed, err := unmarshalSignedMessage(message)
	if err != nil {
		return nil, err
	}
	original, err := rsa.DecryptOAEP(r.hashType.New(), r.rng, privateKey, signed.CipherText, nil)
	if err != nil {
		return nil, err
	}
	digest := r.hashFunc(original)
	err = rsa.VerifyPSS(senderPublicKey, r.hashType, digest, signed.Signature, nil)
	if err != nil {
		return nil, err
	}
	return original, nil
}
