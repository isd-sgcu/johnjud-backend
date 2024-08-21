package jwt

import (
	"fmt"
	"time"

	_jwt "github.com/golang-jwt/jwt/v4"

	"github.com/isd-sgcu/johnjud-gateway/config"
	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/auth/token/strategy"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	"github.com/pkg/errors"
)

type Service interface {
	SignAuth(userId string, role constant.Role, authSessionId string) (string, error)
	VerifyAuth(token string) (*_jwt.Token, error)
	GetConfig() *config.Jwt
}

type serviceImpl struct {
	config   config.Jwt
	strategy strategy.JwtStrategy
	jwtUtil  IJwtUtil
}

func NewService(config config.Jwt, strategy strategy.JwtStrategy, jwtUtil IJwtUtil) Service {
	return &serviceImpl{config: config, strategy: strategy, jwtUtil: jwtUtil}
}

func (s *serviceImpl) SignAuth(userId string, role constant.Role, authSessionId string) (string, error) {
	payloads := dto.AuthPayload{
		RegisteredClaims: _jwt.RegisteredClaims{
			Issuer:    s.config.Issuer,
			ExpiresAt: s.jwtUtil.GetNumericDate(time.Now().Add(time.Second * time.Duration(s.config.ExpiresIn))),
			IssuedAt:  s.jwtUtil.GetNumericDate(time.Now()),
		},
		UserID:        userId,
		AuthSessionID: authSessionId,
	}

	token := s.jwtUtil.GenerateJwtToken(_jwt.SigningMethodHS256, payloads)

	tokenStr, err := s.jwtUtil.SignedTokenString(token, s.config.Secret)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error while signing the token due to: %s", err.Error()))
	}

	return tokenStr, nil
}

func (s *serviceImpl) VerifyAuth(token string) (*_jwt.Token, error) {
	return s.jwtUtil.ParseToken(token, s.strategy.AuthDecode)
}

func (s *serviceImpl) GetConfig() *config.Jwt {
	return &s.config
}
