package image

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) FindByPetId(petId string) (result []*proto.Image, err *dto.ResponseErr) {
	args := s.Called(petId)

	if args.Get(0) != nil {
		result = args.Get(0).([]*proto.Image)
	}

	if args.Get(1) != nil {
		err = args.Get(0).(*dto.ResponseErr)
	}

	return
}
