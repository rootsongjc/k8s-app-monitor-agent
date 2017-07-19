package util

import (
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestGCM(t *testing.T) {
	assert := assert.New(t)
	key := Crypto.CreateKey(32)
	plaintext := []byte("Mary Jane Hawkins")

	nonce, ciphertext, tag, err := GCM.Encrypt(key, plaintext)
	assert.Nil(err)

	decryptedPlaintext, err := GCM.Decrypt(key, nonce, ciphertext, tag)
	assert.Nil(err)
	assert.Equal(string(plaintext), string(decryptedPlaintext))
}
