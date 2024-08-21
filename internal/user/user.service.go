package user

import (
	"errors"
	"net/http"

	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	"github.com/isd-sgcu/johnjud-gateway/internal/model"
	"github.com/isd-sgcu/johnjud-gateway/internal/utils"
	"gorm.io/gorm"
)

type Service interface {
	FindOne(id string) (*dto.User, *dto.ResponseErr)
	Update(id string, request *dto.UpdateUserRequest) (*dto.User, *dto.ResponseErr)
	Delete(id string) (*dto.DeleteUserResponse, *dto.ResponseErr)
}

type serviceImpl struct {
	repo       Repository
	bcryptUtil utils.IBcryptUtil
}

func NewService(repo Repository, bcryptUtil utils.IBcryptUtil) Service {
	return &serviceImpl{repo: repo, bcryptUtil: bcryptUtil}
}

func (s *serviceImpl) FindOne(id string) (*dto.User, *dto.ResponseErr) {
	raw := model.User{}

	err := s.repo.FindById(id, &raw)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    constant.UserNotFoundErrorMessage,
			}
		}
		return nil, &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    "Find user failed",
		}
	}

	return RawToDto(&raw), nil
}

func (s *serviceImpl) Update(id string, request *dto.UpdateUserRequest) (*dto.User, *dto.ResponseErr) {
	hashPassword, err := s.bcryptUtil.GenerateHashedPassword(request.Password)
	if err != nil {
		return nil, &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    constant.InternalServerErrorMessage,
		}
	}

	updateUser := &model.User{
		Email:     request.Email,
		Password:  hashPassword,
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
	}

	err = s.repo.Update(id, updateUser)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusConflict,
				Message:    constant.DuplicateEmailErrorMessage,
			}
		}
		return nil, &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    "Update user failed",
		}
	}

	return RawToDto(updateUser), nil
}

func (s *serviceImpl) Delete(id string) (*dto.DeleteUserResponse, *dto.ResponseErr) {
	err := s.repo.Delete(id)
	if err != nil {
		return nil, &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    "Delete user failed",
		}
	}

	return &dto.DeleteUserResponse{Success: true}, nil
}

func RawToDto(in *model.User) *dto.User {
	return &dto.User{
		Id:        in.ID.String(),
		Email:     in.Email,
		Firstname: in.Firstname,
		Lastname:  in.Lastname,
	}
}
