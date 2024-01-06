package pet

import (
	"context"

	petProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type PetClientMock struct {
	mock.Mock
}

func (c *PetClientMock) AdoptPet(_ context.Context, in *petProto.AdoptPetRequest, _ ...grpc.CallOption) (res *petProto.AdoptPetResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*petProto.AdoptPetResponse)
	}
	return res, args.Error(1)
}

func (c *PetClientMock) FindAll(_ context.Context, in *petProto.FindAllPetRequest, _ ...grpc.CallOption) (res *petProto.FindAllPetResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*petProto.FindAllPetResponse)
	}
	return res, args.Error(1)
}

func (c *PetClientMock) FindOne(_ context.Context, in *petProto.FindOnePetRequest, _ ...grpc.CallOption) (res *petProto.FindOnePetResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*petProto.FindOnePetResponse)
	}
	return res, args.Error(1)
}

func (c *PetClientMock) Create(_ context.Context, in *petProto.CreatePetRequest, _ ...grpc.CallOption) (res *petProto.CreatePetResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*petProto.CreatePetResponse)
	}

	return res, args.Error(1)
}

func (c *PetClientMock) Update(_ context.Context, in *petProto.UpdatePetRequest, _ ...grpc.CallOption) (res *petProto.UpdatePetResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*petProto.UpdatePetResponse)
	}

	return res, args.Error(1)
}

func (c *PetClientMock) ChangeView(_ context.Context, in *petProto.ChangeViewPetRequest, _ ...grpc.CallOption) (res *petProto.ChangeViewPetResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*petProto.ChangeViewPetResponse)
	}

	return res, args.Error(1)
}

func (c *PetClientMock) Delete(_ context.Context, in *petProto.DeletePetRequest, _ ...grpc.CallOption) (res *petProto.DeletePetResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*petProto.DeletePetResponse)
	}

	return res, args.Error(1)
}
