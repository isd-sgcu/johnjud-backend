package user

import (
	"github.com/isd-sgcu/johnjud-gateway/internal/model"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) FindAll(user *[]*model.User) error {
	args := m.Called(user)
	if args.Get(0) != nil {
		*user = *args.Get(0).(*[]*model.User)
		return nil
	}

	return args.Error(1)
}

func (m *UserRepositoryMock) FindById(id string, user *model.User) error {
	args := m.Called(id, user)
	if args.Get(0) != nil {
		*user = *args.Get(0).(*model.User)
		return nil
	}

	return args.Error(1)
}

func (m *UserRepositoryMock) FindByEmail(email string, user *model.User) error {
	args := m.Called(email, user)
	if args.Get(0) != nil {
		*user = *args.Get(0).(*model.User)
		return nil
	}

	return args.Error(1)
}

func (m *UserRepositoryMock) Create(user *model.User) error {
	args := m.Called(user)
	if args.Get(0) != nil {
		*user = *args.Get(0).(*model.User)
		return nil
	}

	return args.Error(1)
}

func (m *UserRepositoryMock) Update(id string, user *model.User) error {
	args := m.Called(id, user)
	if args.Get(0) != nil {
		*user = *args.Get(0).(*model.User)
		return nil
	}

	return args.Error(1)
}

func (m *UserRepositoryMock) Delete(id string) error {
	args := m.Called(id)

	return args.Error(0)
}
