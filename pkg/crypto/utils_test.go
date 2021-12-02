package crypto

import "testing"

func TestRetrieveKey(t *testing.T) {
	_, pub := GenerateKeyPair()
	enc, err := EncodeKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("len enc =", len(enc))
	t.Log(enc)
	dec, err := DecodeKey(enc)

	if pub.N.Cmp(dec.N) != 0 || pub.E != dec.E {
		t.Fail()
	}
}
