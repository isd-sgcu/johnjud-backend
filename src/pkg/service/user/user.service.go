package user

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
)

type Service interface {
	FindOne(string) (*dto.User, *dto.ResponseErr)
	Update(string, *dto.UpdateUserRequest) (*dto.User, *dto.ResponseErr)
	Delete(string) (*dto.DeleteUserResponse, *dto.ResponseErr)
}
