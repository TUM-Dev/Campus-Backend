// Package ios_crypto provides functions for encrypting and decrypting strings using AES-256-GCM.
package ios_crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

type EncryptedString string

func (e *EncryptedString) String() string {
	return string(*e)
}

func SymmetricEncrypt(plaintext string, key string) (*EncryptedString, error) {
	bytesKey := []byte(key)
	bytesPlaintext := []byte(plaintext)

	c, err := aes.NewCipher(bytesKey)

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, bytesPlaintext, nil)

	encryptedString := EncryptedString(base64.StdEncoding.EncodeToString(ciphertext))

	return &encryptedString, nil
}

func SymmetricDecrypt(encryptedString EncryptedString, key string) (string, error) {
	bytesKey := []byte(key)
	bytesEncryptedString, err := base64.StdEncoding.DecodeString(string(encryptedString))

	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(bytesKey)

	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := bytesEncryptedString[:nonceSize], bytesEncryptedString[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
