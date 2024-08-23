package utils

import "github.com/stretchr/testify/mock"

type BcryptUtilMock struct {
	mock.Mock
}

func (m *BcryptUtilMock) GenerateHashedPassword(password string) (string, error) {
	args := m.Called(password)
	if args.Get(0) != "" {
		return args.Get(0).(string), nil
	}

	return "", args.Error(1)
}

func (m *BcryptUtilMock) CompareHashedPassword(hashedPassword string, plainPassword string) error {
	args := m.Called(hashedPassword, plainPassword)
	if args.Get(0) == nil {
		return nil
	}

	return args.Error(0)
}
