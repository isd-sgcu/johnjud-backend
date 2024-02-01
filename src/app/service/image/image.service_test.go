package image

import (
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	utils "github.com/isd-sgcu/johnjud-gateway/src/app/utils/image"
	imageMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/client/image"
	imageProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ImageServiceTest struct {
	suite.Suite
	Images []*imageProto.Image
	PetId  uuid.UUID
}

func TestImageService(t *testing.T) {
	suite.Run(t, new(ImageServiceTest))
}

func (t *ImageServiceTest) SetupTest() {
	t.PetId = uuid.New()
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
