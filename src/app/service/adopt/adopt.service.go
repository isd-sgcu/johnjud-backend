package adopt

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/adopt/v1"
)

type Service struct {
	client proto.AdoptServiceClient
}

func NewService(client proto.AdoptServiceClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) FindAll() ([]*proto.Adopt, *dto.ResponseErr) {
	return nil, nil
}

func (s *Service) Create(in *dto.AdoptDto) (*proto.Adopt, *dto.ResponseErr) {
	return nil, nil
}

func (s *Service) Delete(id string) (bool, *dto.ResponseErr) {
	return false, nil
}
