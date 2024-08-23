package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
	"time"
)

type JwtUtilMock struct {
	mock.Mock
}

func (m *JwtUtilMock) GenerateJwtToken(method jwt.SigningMethod, payloads jwt.Claims) *jwt.Token {
	args := m.Called(method, payloads)
	return args.Get(0).(*jwt.Token)
}

func (m *JwtUtilMock) GetNumericDate(time time.Time) *jwt.NumericDate {
	args := m.Called(time)
	return args.Get(0).(*jwt.NumericDate)
}

func (m *JwtUtilMock) SignedTokenString(token *jwt.Token, secret string) (string, error) {
	args := m.Called(token, secret)
	if args.Get(0) != "" {
		return args.Get(0).(string), nil
	}

	return "", args.Error(1)
}

func (m *JwtUtilMock) ParseToken(tokenStr string, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
	args := m.Called(tokenStr, keyFunc)
	if args.Get(0) != nil {
		return args.Get(0).(*jwt.Token), nil
	}

	return nil, args.Error(1)
}
