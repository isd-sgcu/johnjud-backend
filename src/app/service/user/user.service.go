package user

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
)

type Service struct {
	client proto.UserServiceClient
}

func NewService(client proto.UserServiceClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) FindOne(id string) (*proto.User, *dto.ResponseErr) {
	return nil, nil
}

func (s *Service) Update(id string, in *dto.UpdateUserDto) (*proto.User, *dto.ResponseErr) {
	return nil, nil
}

func (s *Service) Delete(id string) (bool, *dto.ResponseErr) {
	return false, nil
}
