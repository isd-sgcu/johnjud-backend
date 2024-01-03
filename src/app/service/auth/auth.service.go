package auth

import (
	"context"
	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	auth_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/auth/v1"
	"github.com/rs/zerolog/log"
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

	resp, err := s.client.SignUp(ctx, &auth_proto.SignUpRequest{
		FirstName: request.Firstname,
		LastName:  request.Lastname,
		Email:     request.Email,
		Password:  request.Password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		log.Error().
			Str("service", "auth").
			Str("action", "SignUp").
			Str("email", request.Email).
			Msg(st.Message())
		if ok {
			switch st.Code() {
			case codes.AlreadyExists:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusConflict,
					Message:    constant.DuplicateEmailMessage,
					Data:       nil,
				}
			case codes.Internal:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusInternalServerError,
					Message:    constant.InternalErrorMessage,
					Data:       nil,
				}
			default:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    constant.UnavailableServiceMessage,
					Data:       nil,
				}
			}
		}

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    constant.InternalErrorMessage,
			Data:       nil,
		}
	}

	log.Info().
		Str("service", "auth").
		Str("action", "SignUp").
		Str("email", request.Email).
		Msg("sign up successfully")
	return &dto.SignupResponse{
		Id:        resp.Id,
		Email:     resp.Email,
		Firstname: resp.FirstName,
		Lastname:  resp.LastName,
	}, nil
}

func (s *Service) SignIn(signIn *dto.SignInRequest) (*dto.Credential, *dto.ResponseErr) {
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

	return &dto.Credential{
		AccessToken:  resp.Credential.AccessToken,
		RefreshToken: resp.Credential.RefreshToken,
		ExpiresIn:    int(resp.Credential.ExpiresIn),
	}, nil
}

func (s *Service) SignOut(token string) (*dto.SignOutResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.client.SignOut(ctx, &auth_proto.SignOutRequest{
		Token: token,
	})
	if err != nil {
		st, _ := status.FromError(err)
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

	return &dto.SignOutResponse{
		IsSuccess: response.IsSuccess,
	}, nil
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
