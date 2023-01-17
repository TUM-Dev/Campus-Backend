// Package ios_apns_jwt handles the generation and validation of the JWT token for the APNs service
package ios_apns_jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// TokenTimeout for the token in seconds
	TokenTimeout = 3000
)

var (
	ErrorAuthKeyNotPem   = errors.New("failed to parse token: AuthKey must be a valid .p8 PEM file")
	ErrorAuthKeyNotECDSA = errors.New("failed to parse token: AuthKey must be of type ecdsa.PrivateKey")
	ErrorAuthKeyNil      = errors.New("failed to parse token: AuthKey was nil")
	APNsKeyId            = os.Getenv("APNS_KEY_ID")
	APNsTeamId           = os.Getenv("APNS_TEAM_ID")
	APNsP8FilePath       = os.Getenv("APNS_P8_FILE_PATH")
)

type Token struct {
	sync.Mutex
	EncryptionKey *ecdsa.PrivateKey
	KeyId         string
	TeamId        string
	IssuedAt      int64
	Bearer        string
}

func NewToken() (*Token, error) {
	encryptionKey, err := APNsEncryptionKeyFromFile()

	if err != nil {
		return nil, err
	}

	token := Token{
		EncryptionKey: encryptionKey,
		KeyId:         APNsKeyId,
		TeamId:        APNsTeamId,
	}

	_, err = token.Generate()

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// APNsEncryptionKeyFromFile reads the APNs encryption key from the file system
// and returns it as an ecdsa.PrivateKey
// The file location is defined by the APNS_P8_FILE_PATH environment variable
func APNsEncryptionKeyFromFile() (*ecdsa.PrivateKey, error) {
	path, err := filepath.Abs(APNsP8FilePath)

	if err != nil {
		log.Error("No valid path to AuthKey")

		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error("Could not read APNs encryption key from file")

		return nil, err
	}

	block, _ := pem.Decode(data)

	if block == nil {
		log.Error("Could not decode APNs encryption key from file")

		return nil, ErrorAuthKeyNotPem
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)

	if err != nil {
		log.Error("Could not parse APNs encryption key from file")

		return nil, err
	}

	if ecdsaKey, ok := key.(*ecdsa.PrivateKey); ok {
		return ecdsaKey, nil
	}

	return nil, ErrorAuthKeyNotECDSA
}

func (t *Token) GenerateNewTokenIfExpired() (bearer string) {
	t.Lock()
	defer t.Unlock()

	if t.IsExpired() {
		t.Generate()
	}

	return t.Bearer
}

func (t *Token) IsExpired() bool {
	return time.Now().Unix() >= (t.IssuedAt + TokenTimeout)
}

func (t *Token) Generate() (bool, error) {
	if t.EncryptionKey == nil {
		return false, ErrorAuthKeyNil
	}

	issuedAt := time.Now().Unix()

	jwtToken := &jwt.Token{
		Header: map[string]interface{}{
			"alg": "ES256",
			"kid": t.KeyId,
		},
		Claims: jwt.MapClaims{
			"iss": t.TeamId,
			"iat": issuedAt,
		},
		Method: jwt.SigningMethodES256,
	}

	token, err := jwtToken.SignedString(t.EncryptionKey)

	if err != nil {
		return false, err
	}

	t.IssuedAt = issuedAt
	t.Bearer = token

	return true, nil
}
