package image

import (
	"context"

	imageProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type ImageClientMock struct {
	mock.Mock
}

func (c *ImageClientMock) Upload(_ context.Context, in *imageProto.UploadImageRequest, _ ...grpc.CallOption) (res *imageProto.UploadImageResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*imageProto.UploadImageResponse)
	}
	return res, args.Error(1)
}

func (c *ImageClientMock) FindByPetId(_ context.Context, in *imageProto.FindImageByPetIdRequest, _ ...grpc.CallOption) (res *imageProto.FindImageByPetIdResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*imageProto.FindImageByPetIdResponse)
	}
	return res, args.Error(1)
}

func (c *ImageClientMock) AssignPet(_ context.Context, in *imageProto.AssignPetRequest, _ ...grpc.CallOption) (res *imageProto.AssignPetResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*imageProto.AssignPetResponse)
	}
	return res, args.Error(1)
}

func (c *ImageClientMock) Delete(_ context.Context, in *imageProto.DeleteImageRequest, _ ...grpc.CallOption) (res *imageProto.DeleteImageResponse, err error) {
	args := c.Called(in)

	if args.Get(0) != nil {
		res = args.Get(0).(*imageProto.DeleteImageResponse)
	}
	return res, args.Error(1)
}
