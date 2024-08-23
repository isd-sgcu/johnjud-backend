package pet

import (
	"github.com/isd-sgcu/johnjud-backend/internal/model"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) FindOne(id string, result *model.Pet) error {
	args := r.Called(id, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*model.Pet)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Create(in *model.Pet) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*model.Pet)
	}

	return args.Error(1)
}

func (r *RepositoryMock) FindAll(result *[]*model.Pet, isAdmin bool) error {
	args := r.Called(*result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*[]*model.Pet)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Update(id string, result *model.Pet) error {
	args := r.Called(id, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*model.Pet)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Delete(id string) error {
	args := r.Called(id)
	return args.Error(0)
}
