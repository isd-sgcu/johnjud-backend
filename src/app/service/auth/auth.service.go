package auth

import (
	"context"
	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	auth_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

type Service struct {
	client auth_proto.AuthServiceClient
}

func NewService(client auth_proto.AuthServiceClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Signup(request *dto.SignupRequest) (*dto.SignupResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := s.client.Signup(ctx, &auth_proto.SignupRequest{
		FirstName: request.Firstname,
		LastName:  request.Lastname,
		Email:     request.Email,
		Password:  request.Password,
	})
	if err != nil {
		st, _ := status.FromError(err)
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

	return &dto.SignupResponse{
		Id:        resp.Id,
		Email:     resp.Email,
		Firstname: resp.FirstName,
		Lastname:  resp.LastName,
	}, nil
}

func (s *Service) SignIn(signIn *dto.SignIn) (*auth_proto.Credential, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := s.client.SignIn(ctx, &auth_proto.SignInRequest{
		Email:    signIn.Email,
		Password: signIn.Password,
	})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.PermissionDenied:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusForbidden,
				Message:    constant.IncorrectEmailPasswordMessage,
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

	return resp.Credential, nil
}

func (s *Service) SignOut(token string) (bool, *dto.ResponseErr) {
	return false, nil
}

func (s *Service) Validate(token string) (*dto.TokenPayloadAuth, *dto.ResponseErr) {
	// call authClient.Validate()
	// handle error: unauthorized, internal error and unavailable service
	return nil, nil
}

func (s *Service) RefreshToken(token string) (*auth_proto.Credential, *dto.ResponseErr) {
	// call authClient.Validate()
	// handle error: unauthorized, internal error and unavailable service
	return nil, nil
}
