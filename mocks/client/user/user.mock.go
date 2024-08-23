package user

import (
	"context"

	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type UserClientMock struct {
	mock.Mock
}

func (c *UserClientMock) FindOne(_ context.Context, in *proto.FindOneUserRequest, _ ...grpc.CallOption) (res *proto.FindOneUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.FindOneUserResponse)
	}
	return res, args.Error(1)
}
func (c *UserClientMock) Update(_ context.Context, in *proto.UpdateUserRequest, _ ...grpc.CallOption) (res *proto.UpdateUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.UpdateUserResponse)
	}

	return res, args.Error(1)
}

func (c *UserClientMock) Delete(_ context.Context, in *proto.DeleteUserRequest, _ ...grpc.CallOption) (res *proto.DeleteUserResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*proto.DeleteUserResponse)
	}

	return res, args.Error(1)
}
