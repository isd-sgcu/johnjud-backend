package like

import (
	"context"

	imageProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type LikeClientMock struct {
	mock.Mock
}

func (c *LikeClientMock) FindByPetId(_ context.Context, in *imageProto.FindImageByPetIdRequest, _ ...grpc.CallOption) (res *imageProto.FindImageByPetIdResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*imageProto.FindImageByPetIdResponse)
	}
	return res, args.Error(1)
}
func (c *LikeClientMock) Create(_ context.Context, in *imageProto.UploadImageRequest, _ ...grpc.CallOption) (res *imageProto.UploadImageResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*imageProto.UploadImageResponse)
	}
	return res, args.Error(1)
}
func (c *LikeClientMock) Delete(_ context.Context, in *imageProto.DeleteImageRequest, _ ...grpc.CallOption) (res *imageProto.DeleteImageResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*imageProto.DeleteImageResponse)
	}
	return res, args.Error(1)
}
