package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type IJwtUtil interface {
	GenerateJwtToken(method jwt.SigningMethod, payloads jwt.Claims) *jwt.Token
	GetNumericDate(time time.Time) *jwt.NumericDate
	SignedTokenString(token *jwt.Token, secret string) (string, error)
	ParseToken(tokenStr string, keyFunc jwt.Keyfunc) (*jwt.Token, error)
}

type jwtUtilImpl struct{}

func NewJwtUtil() IJwtUtil {
	return &jwtUtilImpl{}
}

func (u *jwtUtilImpl) GenerateJwtToken(method jwt.SigningMethod, payloads jwt.Claims) *jwt.Token {
	return jwt.NewWithClaims(method, payloads)
}

func (u *jwtUtilImpl) GetNumericDate(time time.Time) *jwt.NumericDate {
	return jwt.NewNumericDate(time)
}

func (u *jwtUtilImpl) SignedTokenString(token *jwt.Token, secret string) (string, error) {
	return token.SignedString([]byte(secret))
}

func (u *jwtUtilImpl) ParseToken(tokenStr string, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, keyFunc)
}
