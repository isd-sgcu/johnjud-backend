package user

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
)

type Service interface {
	FindOne(string) (*dto.FindOneUserResponse, *dto.ResponseErr)
	Update(string, *dto.UpdateUserRequest) (*dto.UpdateUserResponse, *dto.ResponseErr)
	Delete(string) (*dto.DeleteUserResponse, *dto.ResponseErr)
}
