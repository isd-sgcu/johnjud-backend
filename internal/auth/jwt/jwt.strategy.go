package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/isd-sgcu/johnjud-gateway/internal/auth/token/strategy"
	"github.com/pkg/errors"
)

type JwtStrategy interface {
	AuthDecode(token *jwt.Token) (interface{}, error)
}

type jwtStrategyImpl struct {
	secret string
}

func NewJwtStrategy(secret string) strategy.JwtStrategy {
	return &jwtStrategyImpl{secret: secret}
}

func (s *jwtStrategyImpl) AuthDecode(token *jwt.Token) (interface{}, error) {
	if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
		return nil, errors.New(fmt.Sprintf("invalid token %v\n", token.Header["alg"]))
	}

	return []byte(s.secret), nil
}
