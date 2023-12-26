package auth

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	auth_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/auth/v1"
)

type Service interface {
	Signup(*dto.SignupRequest) (*dto.SignupResponse, *dto.ResponseErr)
	SignIn(*dto.SignIn) (*dto.Credential, *dto.ResponseErr)
	SignOut(string) (bool, *dto.ResponseErr)
	Validate(string) (*dto.TokenPayloadAuth, *dto.ResponseErr)
	RefreshToken(string) (*auth_proto.Credential, *dto.ResponseErr)
}
