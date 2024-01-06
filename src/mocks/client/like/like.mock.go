package like

import (
	"context"

	likeProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/like/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type LikeClientMock struct {
	mock.Mock
}

func (c *LikeClientMock) FindByUserId(_ context.Context, in *likeProto.FindLikeByUserIdRequest, _ ...grpc.CallOption) (res *likeProto.FindLikeByUserIdResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*likeProto.FindLikeByUserIdResponse)
	}
	return res, args.Error(1)
}
func (c *LikeClientMock) Create(_ context.Context, in *likeProto.CreateLikeRequest, _ ...grpc.CallOption) (res *likeProto.CreateLikeResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*likeProto.CreateLikeResponse)
	}
	return res, args.Error(1)
}
func (c *LikeClientMock) Delete(_ context.Context, in *likeProto.DeleteLikeRequest, _ ...grpc.CallOption) (res *likeProto.DeleteLikeResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*likeProto.DeleteLikeResponse)
	}
	return res, args.Error(1)
}
