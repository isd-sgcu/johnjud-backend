package user

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	user_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
)

type Service interface {
	FindOne(string) (*dto.FindOneUserResponse, *dto.ResponseErr)
	Update(string, *dto.UpdateUserRequest) (*user_proto.User, *dto.ResponseErr)
}
