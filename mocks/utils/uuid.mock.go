package utils

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UuidUtilMock struct {
	mock.Mock
}

func (m *UuidUtilMock) GetNewUUID() *uuid.UUID {
	args := m.Called()
	return args.Get(0).(*uuid.UUID)
}
