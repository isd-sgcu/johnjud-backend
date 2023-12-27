package pet

import (
	"context"

	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) FindOne(id string) (result *proto.Pet, err *dto.ResponseErr) {
	args := s.Called(id)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.Pet)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

func (s *ServiceMock) Create(in *dto.PetDto) (result *proto.Pet, err *dto.ResponseErr) {
	args := s.Called(in)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.Pet)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

func (s *ServiceMock) Update(id string, in *proto.Pet) (result *proto.Pet, err *dto.ResponseErr) {
	args := s.Called(id, in)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.Pet)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

func (s *ServiceMock) Delete(in string) (result *proto.Pet, err *dto.ResponseErr) {
	args := s.Called(in)

	if args.Get(0) != nil {
		result = args.Get(0).(*proto.Pet)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) FindOne(ctx context.Context, in *proto.FindOnePetRequest, opts ...grpc.CallOption) (res *proto.FindOnePetResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindOnePetResponse)
	}
	return res, args.Error(1)
}

func (c *ClientMock) Create(ctx context.Context, in *proto.CreatePetRequest, opts ...grpc.CallOption) (res *proto.CreatePetResponse, err error) {
	args := c.Called(in.Pet)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.CreatePetResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Update(ctx context.Context, in *proto.UpdatePetRequest, opts ...grpc.CallOption) (res *proto.UpdatePetResponse, err error) {
	args := c.Called(in.Pet)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.UpdatePetResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) ChangeView(ctx context.Context, in *proto.ChangeViewPetRequest, opts ...grpc.CallOption) (res *proto.ChangeViewPetResponse, err error) {
	args := c.Called(in.Id)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.ChangeViewPetResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Delete(ctx context.Context, in *proto.DeletePetRequest, opts ...grpc.CallOption) (res *proto.DeletePetResponse, err error) {
	args := c.Called(in.Id)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.DeletePetResponse)
	}

	return res, args.Error(1)
}

type ContextMock struct {
	mock.Mock
	V      interface{}
	Status int
}

func (c *ContextMock) JSON(status int, v interface{}) {
	c.V = v
	c.Status = status
}

func (c *ContextMock) Bind(v interface{}) error {
	args := c.Called(v)

	if args.Get(0) != nil {
		switch v.(type) {
		case *dto.PetDto:
			*v.(*dto.PetDto) = *args.Get(0).(*dto.PetDto)
		}
	}

	return args.Error(1)
}

func (c *ContextMock) ID() (string, error) {
	args := c.Called()
	return args.String(0), args.Error(1)
}

func (c *ContextMock) Host() string {
	args := c.Called()
	return args.String(0)
}

func (c *ContextMock) PetID() string {
	args := c.Called()
	return args.String(0)
}

func (c *ContextMock) Query(key string) string {
	args := c.Called(key)
	return args.String(0)
}
