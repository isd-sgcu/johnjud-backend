package mock_image

import (
	dto "github.com/isd-sgcu/johnjud-backend/internal/dto"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (c *ServiceMock) FindAll() (res []*dto.ImageResponse, err *dto.ResponseErr) {
	args := c.Called()

	if args.Get(0) != nil {
		res = args.Get(0).([]*dto.ImageResponse)
	}

	return res, args.Get(1).(*dto.ResponseErr)
}

func (c *ServiceMock) FindByPetId(petID string) ([]*dto.ImageResponse, *dto.ResponseErr) {
	args := c.Called(petID)

	if args.Get(0) != nil {
		res := args.Get(0).([]*dto.ImageResponse)
		return res, nil
	}
	return nil, args.Get(1).(*dto.ResponseErr)
}

func (c *ServiceMock) Upload(request *dto.UploadImageRequest) (*dto.ImageResponse, *dto.ResponseErr) {
	args := c.Called(request)

	if args.Get(0) != nil {
		res := args.Get(0).(*dto.ImageResponse)
		return res, nil
	}
	return nil, args.Get(1).(*dto.ResponseErr)
}

func (c *ServiceMock) Delete(id string) (*dto.DeleteImageResponse, *dto.ResponseErr) {
	args := c.Called(id)

	if args.Get(0) != nil {
		res := args.Get(0).(*dto.DeleteImageResponse)
		return res, nil
	}
	return nil, args.Get(1).(*dto.ResponseErr)
}

func (c *ServiceMock) DeleteByPetId(petID string) (*dto.DeleteImageResponse, *dto.ResponseErr) {
	args := c.Called(petID)

	if args.Get(0) != nil {
		res := args.Get(0).(*dto.DeleteImageResponse)
		return res, nil
	}
	return nil, args.Get(1).(*dto.ResponseErr)
}

func (c *ServiceMock) AssignPet(request *dto.AssignPetRequest) (*dto.AssignPetResponse, *dto.ResponseErr) {
	args := c.Called(request)

	if args.Get(0) != nil {
		res := args.Get(0).(*dto.AssignPetResponse)
		return res, nil
	}
	return nil, args.Get(1).(*dto.ResponseErr)
}
