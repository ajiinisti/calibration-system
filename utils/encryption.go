package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func EncryptUUID(uuid string, secretKey string) (string, error) {
	// Creating block of algorithm
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("cipher err: %v", err.Error())
	}

	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("cipher GCM err: %v", err.Error())
	}

	// Generating random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce  err: %v", err.Error())
	}

	// Decrypt file
	cipherText := gcm.Seal(nonce, nonce, []byte(uuid), nil)

	return string(cipherText), nil

}

func DecryptUUID(token string, secretKey string) (string, error) {
	// Creating block of algorithm
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("cipher err: %v", err.Error())
	}

	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("cipher GCM err: %v", err.Error())
	}

	// Ensure that the ciphertext is not empty
	if len(token) <= gcm.NonceSize() {
		return "", fmt.Errorf("invalid ciphertext")
	}

	// Deattached nonce and decrypt
	cipherText := []byte(token)
	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]

	// Decrypt and authenticate the message
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt file err: %v", err.Error())
	}
	return string(plainText), nil
}
