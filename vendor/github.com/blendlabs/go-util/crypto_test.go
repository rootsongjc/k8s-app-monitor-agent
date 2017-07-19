package util

import (
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestCryptoEncryptDecrypt(t *testing.T) {
	assert := assert.New(t)
	key := Crypto.CreateKey(32)
	plaintext := "Mary Jane Hawkins"

	ciphertext, err := Crypto.Encrypt(key, []byte(plaintext))
	assert.Nil(err)

	decryptedPlaintext, err := Crypto.Decrypt(key, ciphertext)
	assert.Nil(err)
	assert.Equal(plaintext, string(decryptedPlaintext))
}

func TestCryptoHash(t *testing.T) {
	assert := assert.New(t)
	key := Crypto.CreateKey(128)
	plaintext := "123-12-1234"
	assert.Equal(
		Crypto.Hash(key, []byte(plaintext)),
		Crypto.Hash(key, []byte(plaintext)),
	)
}
