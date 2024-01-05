package user

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	client proto.UserServiceClient
}

func NewService(client proto.UserServiceClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) FindOne(id string) (*dto.FindOneUserResponse, *dto.ResponseErr) {
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

	return &dto.FindOneUserResponse{
		Id:        response.User.Id,
		Firstname: response.User.Firstname,
		Lastname:  response.User.Lastname,
		Email:     response.User.Email,
	}, nil
}

func (s *Service) Update(id string, in *dto.UpdateUserRequest) (*dto.UpdateUserResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.client.Update(ctx, &proto.UpdateUserRequest{
		Id:        id,
		Firstname: in.Firstname,
		Lastname:  in.Lastname,
	})
	if err != nil {
		st, _ := status.FromError(err)
		log.Error().
			Err(err).
			Str("service", "user").
			Str("module", "update").
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

	return &dto.UpdateUserResponse{
		Id:        response.User.Id,
		Firstname: response.User.Firstname,
		Lastname:  response.User.Lastname,
		Email:     response.User.Email,
	}, nil
}
