package auth

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	auth_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/auth/v1"
	user_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
)

type Service struct {
	client     auth_proto.AuthServiceClient
	userClient user_proto.UserServiceClient
}

func NewService(client auth_proto.AuthServiceClient, userClient user_proto.UserServiceClient) *Service {
	return &Service{
		client:     client,
		userClient: userClient,
	}
}

func (s *Service) Signup(signup *dto.Signup) (*auth_proto.Credential, *dto.ResponseErr) {
	return nil, nil
}

func (s *Service) Signin(signin *dto.Signin) (*auth_proto.Credential, *dto.ResponseErr) {
	return nil, nil
}

func (s *Service) Signout(token string) (bool, *dto.ResponseErr) {
	return false, nil
}

func (s *Service) Validate(token string) (*dto.TokenPayloadAuth, *dto.ResponseErr) {
	return nil, nil
}

func (s *Service) RefreshToken(token string) (*auth_proto.Credential, *dto.ResponseErr) {
	return nil, nil
}
