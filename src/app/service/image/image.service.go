package image

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
)

type Service struct {
	client proto.ImageServiceClient
}

func NewService(client proto.ImageServiceClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) FindByPetId(string) ([]*proto.Image, *dto.ResponseErr) {
	return nil, nil
}

func (s *Service) Upload(in *dto.ImageDto) (*proto.Image, *dto.ResponseErr) {
	return nil, nil
}

func (s *Service) Delete(id string) (bool, *dto.ResponseErr) {
	return false, nil
}
