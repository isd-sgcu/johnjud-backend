package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	authProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/auth/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	client authProto.AuthServiceClient
}

func NewService(client authProto.AuthServiceClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Signup(request *dto.SignupRequest) (*dto.SignupResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := s.client.SignUp(ctx, &authProto.SignUpRequest{
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
		if !ok {
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    constant.InternalErrorMessage,
				Data:       nil,
			}
		}
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

func (s *Service) SignIn(request *dto.SignInRequest) (*dto.Credential, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := s.client.SignIn(ctx, &authProto.SignInRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		log.Error().
			Str("service", "auth").
			Str("action", "SignIn").
			Str("email", request.Email).
			Msg(st.Message())
		if !ok {
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    constant.InternalErrorMessage,
				Data:       nil,
			}
		}
		switch st.Code() {
		case codes.PermissionDenied:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusForbidden,
				Message:    constant.IncorrectEmailPasswordMessage,
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

	log.Info().
		Str("service", "auth").
		Str("action", "SignIn").
		Str("email", request.Email).
		Msg("sign in successfully")
	return &dto.Credential{
		AccessToken:  resp.Credential.AccessToken,
		RefreshToken: resp.Credential.RefreshToken,
		ExpiresIn:    int(resp.Credential.ExpiresIn),
	}, nil
}

func (s *Service) SignOut(token string) (*dto.SignOutResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.client.SignOut(ctx, &authProto.SignOutRequest{
		Token: token,
	})
	if err != nil {
		st, ok := status.FromError(err)
		log.Error().
			Str("service", "auth").
			Str("action", "SignOut").
			Str("token", token).
			Msg(st.Message())
		if !ok {
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    constant.InternalErrorMessage,
				Data:       nil,
			}
		}
		switch st.Code() {
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

	log.Info().
		Str("service", "auth").
		Str("action", "SignOut").
		Str("token", token).
		Msg("sign out successfully")
	return &dto.SignOutResponse{
		IsSuccess: response.IsSuccess,
	}, nil
}

func (s *Service) Validate(token string) (*dto.TokenPayloadAuth, *dto.ResponseErr) {
	// call authClient.Validate()
	// handle error: unauthorized, internal error and unavailable service
	return nil, nil
}

func (s *Service) RefreshToken(request *dto.RefreshTokenRequest) (*dto.Credential, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.client.RefreshToken(ctx, &auth_proto.RefreshTokenRequest{
		RefreshToken: request.RefreshToken,
	})
	if err != nil {
		st, ok := status.FromError(err)
		log.Error().
			Str("service", "auth").
			Str("action", "RefreshToken").
			Str("token", request.RefreshToken).
			Msg(st.Message())
		if !ok {
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    constant.InternalErrorMessage,
				Data:       nil,
			}
		}
		switch st.Code() {
		case codes.InvalidArgument:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusBadRequest,
				Message:    constant.InvalidTokenMessage,
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

	log.Info().
		Str("service", "auth").
		Str("action", "RefreshToken").
		Str("token", request.RefreshToken).
		Msg("Refresh token successfully")
	return &dto.Credential{
		AccessToken:  response.Credential.AccessToken,
		RefreshToken: response.Credential.RefreshToken,
		ExpiresIn:    int(response.Credential.ExpiresIn),
	}, nil
}
