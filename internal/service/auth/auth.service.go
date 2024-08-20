package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.client.Validate(ctx, &authProto.ValidateRequest{
		Token: token,
	})
	if err != nil {
		st, ok := status.FromError(err)
		log.Error().
			Str("service", "auth").
			Str("action", "Validate").
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
		case codes.Unauthenticated:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusUnauthorized,
				Message:    constant.UnauthorizedMessage,
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
		Str("action", "Validate").
		Str("token", token).
		Msg("validate successfully")
	return &dto.TokenPayloadAuth{
		UserId: response.UserId,
		Role:   response.Role,
	}, nil
}

func (s *Service) RefreshToken(request *dto.RefreshTokenRequest) (*dto.Credential, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.client.RefreshToken(ctx, &authProto.RefreshTokenRequest{
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

func (s *Service) ForgotPassword(request *dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.client.ForgotPassword(ctx, &authProto.ForgotPasswordRequest{
		Email: request.Email,
	})
	if err != nil {
		st, ok := status.FromError(err)
		log.Error().
			Str("service", "auth").
			Str("action", "ForgotPassword").
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
		case codes.NotFound:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    constant.UserNotFoundMessage,
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
		Str("action", "ForgotPassword").
		Str("email", request.Email).
		Msg("Forgot password successfully")
	return &dto.ForgotPasswordResponse{
		IsSuccess: true,
	}, nil
}

func (s *Service) ResetPassword(request *dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.client.ResetPassword(ctx, &authProto.ResetPasswordRequest{
		Token:    request.Token,
		Password: request.Password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		log.Error().
			Str("service", "auth").
			Str("action", "ResetPassword").
			Str("token", request.Token).
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
				Message:    constant.ForbiddenSamePasswordMessage,
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
		Str("action", "ResetPassword").
		Str("token", request.Token).
		Msg("Reset password successfully")
	return &dto.ResetPasswordResponse{
		IsSuccess: response.IsSuccess,
	}, nil
}
