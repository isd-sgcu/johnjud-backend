package user

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	user_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
)

type Service interface {
	FindOne(string) (*user_proto.User, *dto.ResponseErr)
	Update(string, *dto.UpdateUserDto) (*user_proto.User, *dto.ResponseErr)
}
