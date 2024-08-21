package user

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	FindOne(string) (*dto.User, *dto.ResponseErr)
	Update(string, *dto.UpdateUserRequest) (*dto.User, *dto.ResponseErr)
	Delete(string) (*dto.DeleteUserResponse, *dto.ResponseErr)
}

type serviceImpl struct {
	client proto.UserServiceClient
}

func NewService(client proto.UserServiceClient) *serviceImpl {
	return &serviceImpl{
		client: client,
	}
}

func (s *serviceImpl) FindOne(id string) (*dto.User, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.client.FindOne(ctx, &proto.FindOneUserRequest{
		Id: id,
	})
	if err != nil {
		st, _ := status.FromError(err)
		log.Error().
			Err(err).
			Str("service", "user").
			Str("module", "find one").
			Msg(st.Message())
		switch st.Code() {
		case codes.NotFound:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    constant.UserNotFoundMessage,
				Data:       nil,
			}

		case codes.Unavailable:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.UnavailableServiceMessage,
				Data:       nil,
			}
		default:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    constant.InternalErrorMessage,
				Data:       nil,
			}
		}
	}

	return &dto.User{
		Id:        response.User.Id,
		Firstname: response.User.Firstname,
		Lastname:  response.User.Lastname,
		Email:     response.User.Email,
	}, nil
}

func (s *serviceImpl) Update(id string, in *dto.UpdateUserRequest) (*dto.User, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.client.Update(ctx, &proto.UpdateUserRequest{
		Id:        id,
		Firstname: in.Firstname,
		Lastname:  in.Lastname,
		Password:  in.Password,
		Email:     in.Email,
	})
	if err != nil {
		st, _ := status.FromError(err)
		log.Error().
			Err(err).
			Str("service", "user").
			Str("module", "update").
			Msg(st.Message())
		switch st.Code() {
		case codes.AlreadyExists:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusConflict,
				Message:    constant.DuplicateEmailMessage,
				Data:       nil,
			}

		case codes.Unavailable:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.UnavailableServiceMessage,
				Data:       nil,
			}
		default:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    constant.InternalErrorMessage,
				Data:       nil,
			}
		}
	}

	return &dto.User{
		Id:        response.User.Id,
		Firstname: response.User.Firstname,
		Lastname:  response.User.Lastname,
		Email:     response.User.Email,
	}, nil
}

func (s *serviceImpl) Delete(id string) (*dto.DeleteUserResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.client.Delete(ctx, &proto.DeleteUserRequest{
		Id: id,
	})
	if err != nil {
		st, _ := status.FromError(err)
		log.Error().
			Err(err).
			Str("service", "user").
			Str("module", "delete").
			Msg(st.Message())
		switch st.Code() {
		case codes.Unavailable:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    constant.UnavailableServiceMessage,
				Data:       nil,
			}
		default:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    constant.InternalErrorMessage,
				Data:       nil,
			}
		}
	}

	return &dto.DeleteUserResponse{
		Success: response.Success,
	}, nil
}
