package common

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

// mostly from go pkg docs
func Decrypt(data, key string) ([]byte, error) {
	k, _ := hex.DecodeString(key)
	ciphertext, _ := hex.DecodeString(data)

	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}
