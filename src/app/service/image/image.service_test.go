package image

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	imageUtils "github.com/isd-sgcu/johnjud-gateway/src/app/utils/image"
	imageProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/stretchr/testify/suite"
)

type ImageServiceTest struct {
	suite.Suite
	ImagesList          [][]*imageProto.Image
	Images              []*imageProto.Image
	ImagesListResponse  [][]*dto.ImageResponse
	ImagesResponse      []*dto.ImageResponse
	UploadImageProtoReq *imageProto.UploadImageRequest
	UploadImageDtoReq   *dto.UploadImageRequest
	DeleteImageProtoReq *imageProto.DeleteImageRequest

	NotFoundErr           *dto.ResponseErr
	UnavailableServiceErr *dto.ResponseErr
	InvalidArgumentErr    *dto.ResponseErr
	InternalErr           *dto.ResponseErr
}

func TestImageService(t *testing.T) {
	suite.Run(t, new(ImageServiceTest))
}

func (t *ImageServiceTest) SetupTest() {
	imagesList := imageUtils.MockImageList(3)
	t.ImagesList = imagesList
	t.Images = imagesList[0]

	t.UploadImageProtoReq = &imageProto.UploadImageRequest{
		Filename: faker.FirstName(),
		Data:     []byte(faker.Word()),
		PetId:    faker.UUIDDigit(),
	}

	t.UploadImageDtoReq = &dto.UploadImageRequest{
		Filename: faker.FirstName(),
		Data:     []byte(faker.Word()),
		PetId:    faker.UUIDDigit(),
	}

	t.DeleteImageProtoReq = &imageProto.DeleteImageRequest{
		Id:        faker.UUIDDigit(),
		ObjectKey: faker.Word(),
	}

	t.UnavailableServiceErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    constant.UnavailableServiceMessage,
		Data:       nil,
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    constant.UserNotFoundMessage,
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
