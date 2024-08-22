package utils

import "github.com/stretchr/testify/mock"

type RandomUtilMock struct {
	mock.Mock
}

func (m *RandomUtilMock) GenerateRandomString(length int) (string, error) {
	args := m.Called(length)
	if args.Get(0) != "" {
		return args.Get(0).(string), nil
	}

	return "", args.Error(1)
}
