package crypto

import (
	"testing"
)

func TestRetrieveMessage(t *testing.T) {
	priv1, pub1 := GenerateKeyPair()
	priv2, pub2 := GenerateKeyPair()

	msg := "a100characterstring.a100characterstring.a100characterstring.a100characterstring.a100characterstring."

	cipher := NewCipher()

	enc1, err := cipher.Encrypt([]byte(msg), priv1, pub2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("length of enc =", len(enc1))
	dec1, err := cipher.Decrypt(enc1, priv2, pub1)
	if err != nil {
		t.Fatal(err)
	}
	if string(dec1) != msg {
		t.Fail()
	}
}
