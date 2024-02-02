package image

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	utils "github.com/isd-sgcu/johnjud-gateway/src/app/utils/image"
	imageMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/client/image"
	imageProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ImageServiceTest struct {
	suite.Suite
	Image          *imageProto.Image
	Images         []*imageProto.Image
	PetId          uuid.UUID
	UploadProtoReq *imageProto.UploadImageRequest
	UploadDtoReq   *dto.UploadImageRequest

	NotFoundErr           *dto.ResponseErr
	UnavailableServiceErr *dto.ResponseErr
	InternalErr           *dto.ResponseErr
	InvalidArgumentErr    *dto.ResponseErr
}

func TestImageService(t *testing.T) {
	suite.Run(t, new(ImageServiceTest))
}

func (t *ImageServiceTest) SetupTest() {
	t.PetId = uuid.New()
	t.Image = &imageProto.Image{
		Id:        faker.UUIDDigit(),
		PetId:     t.PetId.String(),
		ImageUrl:  faker.URL(),
		ObjectKey: faker.Word(),
	}
	t.Images = []*imageProto.Image{
		{
			Id:        faker.UUIDDigit(),
			PetId:     t.PetId.String(),
			ImageUrl:  faker.URL(),
			ObjectKey: faker.Word(),
		}, {
			Id:        faker.UUIDDigit(),
			PetId:     t.PetId.String(),
			ImageUrl:  faker.URL(),
			ObjectKey: faker.Word(),
		}, {
			Id:        faker.UUIDDigit(),
			PetId:     t.PetId.String(),
			ImageUrl:  faker.URL(),
			ObjectKey: faker.Word(),
		}, {
			Id:        faker.UUIDDigit(),
			PetId:     t.PetId.String(),
			ImageUrl:  faker.URL(),
			ObjectKey: faker.Word(),
		},
	}

	t.UploadProtoReq = &imageProto.UploadImageRequest{
		Filename: faker.Name(),
		Data:     []byte{1, 2, 3, 4, 5},
		PetId:    t.PetId.String(),
	}
	t.UploadDtoReq = &dto.UploadImageRequest{
		Filename: t.UploadProtoReq.Filename,
		File:     t.UploadProtoReq.Data,
		PetId:    t.UploadProtoReq.PetId,
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    constant.PetNotFoundMessage,
		Data:       nil,
	}
	t.UnavailableServiceErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    constant.UnavailableServiceMessage,
		Data:       nil,
	}
	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}
	t.InvalidArgumentErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.InvalidArgumentMessage,
		Data:       nil,
	}
}

func (t *ImageServiceTest) TestFindByPetIdSuccess() {
	protoReq := &imageProto.FindImageByPetIdRequest{
		PetId: t.PetId.String(),
	}
	protoResp := &imageProto.FindImageByPetIdResponse{
		Images: t.Images,
	}

	expected := utils.ProtoToDtoList(t.Images)

	client := imageMock.ImageClientMock{}
	client.On("FindByPetId", protoReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.FindByPetId(t.PetId.String())

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *ImageServiceTest) TestFindByPetIdNotFoundError() {
	protoReq := &imageProto.FindImageByPetIdRequest{
		PetId: t.PetId.String(),
	}

	clientErr := status.Error(codes.NotFound, constant.PetNotFoundMessage)

	expected := t.NotFoundErr

	client := imageMock.ImageClientMock{}
	client.On("FindByPetId", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.FindByPetId(t.PetId.String())

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *ImageServiceTest) TestFindByPetIdUnavailableServiceError() {
	protoReq := &imageProto.FindImageByPetIdRequest{
		PetId: t.PetId.String(),
	}

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := imageMock.ImageClientMock{}
	client.On("FindByPetId", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.FindByPetId(t.PetId.String())

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *ImageServiceTest) TestFindByPetIdInternalError() {
	protoReq := &imageProto.FindImageByPetIdRequest{
		PetId: t.PetId.String(),
	}

	clientErr := status.Error(codes.Internal, constant.InternalErrorMessage)

	expected := t.InternalErr

	client := imageMock.ImageClientMock{}
	client.On("FindByPetId", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.FindByPetId(t.PetId.String())

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *ImageServiceTest) TestUploadSuccess() {
	protoReq := t.UploadProtoReq
	protoResp := &imageProto.UploadImageResponse{
		Image: t.Image,
	}

	expected := utils.ProtoToDto(t.Image)

	client := imageMock.ImageClientMock{}
	client.On("Upload", protoReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.Upload(t.UploadDtoReq)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *ImageServiceTest) TestUploadInvalidArgumentError() {
	protoReq := t.UploadProtoReq

	clientErr := status.Error(codes.InvalidArgument, constant.InvalidArgumentMessage)

	expected := t.InvalidArgumentErr

	client := imageMock.ImageClientMock{}
	client.On("Upload", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Upload(t.UploadDtoReq)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *ImageServiceTest) TestUploadUnavailableServiceError() {
	protoReq := t.UploadProtoReq

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := imageMock.ImageClientMock{}
	client.On("Upload", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Upload(t.UploadDtoReq)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *ImageServiceTest) TestUploadInternalError() {
	protoReq := t.UploadProtoReq

	clientErr := status.Error(codes.Internal, constant.InternalErrorMessage)

	expected := t.InternalErr

	client := imageMock.ImageClientMock{}
	client.On("Upload", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Upload(t.UploadDtoReq)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}
