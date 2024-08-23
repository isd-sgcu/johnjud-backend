package strategy

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
)

type JwtStrategyMock struct {
	mock.Mock
}

func (m *JwtStrategyMock) AuthDecode(token *jwt.Token) (interface{}, error) {
	args := m.Called(token)
	if args.Get(0) != nil {
		return args.Get(0), nil
	}
	return nil, args.Error(1)
}
