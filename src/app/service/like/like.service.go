package like

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/like/v1"
)

type Service struct {
	client proto.LikeServiceClient
}

func NewService(client proto.LikeServiceClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) FindByUserId(userId string) ([]*proto.Like, *dto.ResponseErr) {
	return nil, nil
}

func (s *Service) Create(in *dto.LikeDto) (*proto.Like, *dto.ResponseErr) {
	return nil, nil
}

func (s *Service) Delete(id string) (bool, *dto.ResponseErr) {
	return false, nil
}
