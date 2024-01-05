package user

import (
	"context"

	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	user_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) FindOne(id string) (result *user_proto.User, err *dto.ResponseErr) {
	args := s.Called(id)

	if args.Get(0) != nil {
		result = args.Get(0).(*user_proto.User)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

func (s *ServiceMock) Update(id string, in *dto.UpdateUserRequest) (result *user_proto.User, err *dto.ResponseErr) {
	args := s.Called(id, in)

	if args.Get(0) != nil {
		result = args.Get(0).(*user_proto.User)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(*dto.ResponseErr)
	}

	return
}

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) FindOne(_ context.Context, in *user_proto.FindOneUserRequest, _ ...grpc.CallOption) (res *user_proto.FindOneUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*user_proto.FindOneUserResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Update(_ context.Context, in *user_proto.UpdateUserRequest, _ ...grpc.CallOption) (res *user_proto.UpdateUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*user_proto.UpdateUserResponse)
	}

	return res, args.Error(1)
}

func (c *ClientMock) Delete(_ context.Context, in *user_proto.DeleteUserRequest, _ ...grpc.CallOption) (res *user_proto.DeleteUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*user_proto.DeleteUserResponse)
	}

	return res, args.Error(1)
}
