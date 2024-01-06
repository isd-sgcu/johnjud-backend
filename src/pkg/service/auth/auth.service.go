package auth

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
)

type Service interface {
	Signup(*dto.SignupRequest) (*dto.SignupResponse, *dto.ResponseErr)
	SignIn(*dto.SignInRequest) (*dto.Credential, *dto.ResponseErr)
	SignOut(string) (*dto.SignOutResponse, *dto.ResponseErr)
	Validate(string) (*dto.TokenPayloadAuth, *dto.ResponseErr)
	RefreshToken(*dto.RefreshTokenRequest) (*dto.Credential, *dto.ResponseErr)
}
