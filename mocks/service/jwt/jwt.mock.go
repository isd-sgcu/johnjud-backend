package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/isd-sgcu/johnjud-backend/config"
	"github.com/isd-sgcu/johnjud-backend/constant"
	"github.com/stretchr/testify/mock"
)

type JwtServiceMock struct {
	mock.Mock
}

func (m *JwtServiceMock) SignAuth(userId string, role constant.Role, authSessionId string) (string, error) {
	args := m.Called(userId, role, authSessionId)
	if args.Get(0) != "" {
		return args.Get(0).(string), nil
	}

	return "", args.Error(1)
}

func (m *JwtServiceMock) VerifyAuth(token string) (*jwt.Token, error) {
	args := m.Called(token)
	if args.Get(0) != nil {
		return args.Get(0).(*jwt.Token), nil
	}

	return nil, args.Error(1)
}

func (m *JwtServiceMock) GetConfig() *config.Jwt {
	args := m.Called()
	return args.Get(0).(*config.Jwt)
}
