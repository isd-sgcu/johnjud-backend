package utils

import (
	"crypto/rand"
	"encoding/base64"
)

type RandomUtil interface {
	GenerateRandomString(length int) (string, error)
}

func NewRandomUtil() RandomUtil {
	return &randomUtil{}
}

type randomUtil struct{}

func (u *randomUtil) GenerateRandomString(length int) (string, error) {
	numBytes := (length * 6) / 8

	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	return randomString[:length], nil
}
