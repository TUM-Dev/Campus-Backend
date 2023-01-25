// Package ios_crypto provides functions for encrypting and decrypting strings using AES-256-GCM.
package ios_crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	log "github.com/sirupsen/logrus"
)

type EncryptedString string

func (e *EncryptedString) String() string {
	return string(*e)
}

func AsymmetricEncrypt(plaintext string, publicKey string) (*EncryptedString, error) {
	pub, err := StringToPublicKey(publicKey)

	if err != nil {
		return nil, err
	}

	oaep, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pub,
		[]byte(plaintext),
		nil,
	)

	if err != nil {
		return nil, err
	}

	encryptedString := EncryptedString(base64.StdEncoding.EncodeToString(oaep))

	return &encryptedString, nil
}

func StringToPublicKey(pub string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pub))

	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, errors.New("failed to parse DER encoded public key: " + err.Error())
	}

	if pubKey, ok := key.(*rsa.PublicKey); ok {
		return pubKey, nil
	} else {
		return nil, errors.New("failed to parse DER encoded public key: " + err.Error())
	}
}

func SymmetricEncrypt(plaintext string, key string) (*EncryptedString, error) {
	bytesKey := []byte(key)
	bytesPlaintext := []byte(plaintext)

	c, err := aes.NewCipher(bytesKey)
	if err != nil {
		log.WithError(err).Error("Could not create cipher")
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.WithError(err).Error("Could not create GCM")
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = rand.Read(nonce); err != nil {
		log.WithError(err).Error("Could not read random bytes")
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
