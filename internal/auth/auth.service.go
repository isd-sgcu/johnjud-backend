package auth

import (
	"fmt"
	"net/http"

	"github.com/isd-sgcu/johnjud-gateway/config"
	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/auth/email"
	"github.com/isd-sgcu/johnjud-gateway/internal/auth/token"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	"github.com/isd-sgcu/johnjud-gateway/internal/model"
	"github.com/isd-sgcu/johnjud-gateway/internal/user"
	"github.com/isd-sgcu/johnjud-gateway/internal/utils"
	"github.com/rs/zerolog/log"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Service interface {
	Validate(refreshToken string) (*dto.TokenPayloadAuth, *dto.ResponseErr)
	RefreshToken(request *dto.RefreshTokenRequest) (*dto.Credential, *dto.ResponseErr)
	Signup(request *dto.SignupRequest) (*dto.SignupResponse, *dto.ResponseErr)
	SignIn(request *dto.SignInRequest) (*dto.Credential, *dto.ResponseErr)
	SignOut(accessToken string) (*dto.SignOutResponse, *dto.ResponseErr)
	ForgotPassword(request *dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, *dto.ResponseErr)
	ResetPassword(request *dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, *dto.ResponseErr)
}

type serviceImpl struct {
	authRepo     Repository
	userRepo     user.Repository
	tokenService token.Service
	emailService email.Service
	bcryptUtil   utils.IBcryptUtil
	config       config.Auth
}

func NewService(authRepo Repository, userRepo user.Repository, tokenService token.Service, emailService email.Service, bcryptUtil utils.IBcryptUtil, config config.Auth) Service {
	return &serviceImpl{
		authRepo:     authRepo,
		userRepo:     userRepo,
		tokenService: tokenService,
		emailService: emailService,
		bcryptUtil:   bcryptUtil,
		config:       config,
	}
}

func (s *serviceImpl) Validate(refreshToken string) (*dto.TokenPayloadAuth, *dto.ResponseErr) {
	userCredential, err := s.tokenService.Validate(refreshToken)
	if err != nil {
		return nil, &dto.ResponseErr{
			StatusCode: http.StatusUnauthorized,
			Message:    constant.InvalidTokenErrorMessage,
		}
	}

	return &dto.TokenPayloadAuth{
		UserId: userCredential.UserID,
		Role:   string(userCredential.Role),
	}, nil
}

func (s *serviceImpl) RefreshToken(request *dto.RefreshTokenRequest) (*dto.Credential, *dto.ResponseErr) {
	refreshTokenCache, err := s.tokenService.FindRefreshTokenCache(request.RefreshToken)
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.InvalidArgument:
			return nil, dto.BadRequestError(constant.InvalidTokenErrorMessage)
		default:
			return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
		}
	}
	credential, err := s.tokenService.CreateCredential(refreshTokenCache.UserID, refreshTokenCache.Role, refreshTokenCache.AuthSessionID)
	if err != nil {
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	err = s.tokenService.RemoveRefreshTokenCache(request.RefreshToken)
	if err != nil {
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	return credential, nil
}

func (s *serviceImpl) Signup(request *dto.SignupRequest) (*dto.SignupResponse, *dto.ResponseErr) {
	hashPassword, err := s.bcryptUtil.GenerateHashedPassword(request.Password)
	if err != nil {
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	createUser := &model.User{
		Email:     request.Email,
		Password:  hashPassword,
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
		Role:      constant.USER,
	}
	err = s.userRepo.Create(createUser)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, dto.ConflictError(constant.DuplicateEmailErrorMessage)
		}
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	return &dto.SignupResponse{
		Id:        createUser.ID.String(),
		Firstname: createUser.Firstname,
		Lastname:  createUser.Lastname,
		Email:     createUser.Email,
	}, nil
}

func (s *serviceImpl) SignIn(request *dto.SignInRequest) (*dto.Credential, *dto.ResponseErr) {
	user := &model.User{}
	err := s.userRepo.FindByEmail(request.Email, user)
	if err != nil {
		return nil, dto.UnauthorizedError(constant.IncorrectEmailPasswordErrorMessage)
	}

	err = s.bcryptUtil.CompareHashedPassword(user.Password, request.Password)
	if err != nil {
		return nil, dto.UnauthorizedError(constant.IncorrectEmailPasswordErrorMessage)
	}

	createAuthSession := &model.AuthSession{
		UserID: user.ID,
	}
	err = s.authRepo.Create(createAuthSession)
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "auth").
			Str("module", "signin").
			Msg("Error creating auth session")
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	credential, err := s.tokenService.CreateCredential(user.ID.String(), user.Role, createAuthSession.ID.String())
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "auth").
			Str("module", "signin").
			Msg("Error creating credential")
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	return credential, nil
}

func (s *serviceImpl) SignOut(accessToken string) (*dto.SignOutResponse, *dto.ResponseErr) {
	userCredential, err := s.tokenService.Validate(accessToken)
	if err != nil {
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	err = s.tokenService.RemoveRefreshTokenCache(userCredential.RefreshToken)
	if err != nil {
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	err = s.tokenService.RemoveAccessTokenCache(userCredential.AuthSessionID)
	if err != nil {
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	err = s.authRepo.Delete(userCredential.AuthSessionID)
	if err != nil {
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	return &dto.SignOutResponse{IsSuccess: true}, nil
}

func (s *serviceImpl) ForgotPassword(request *dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, *dto.ResponseErr) {
	user := &model.User{}
	err := s.userRepo.FindByEmail(request.Email, user)
	if err != nil {
		return nil, dto.NotFoundError(constant.UserNotFoundErrorMessage)
	}

	resetPasswordToken, err := s.tokenService.CreateResetPasswordToken(user.ID.String())
	if err != nil {
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	resetPasswordURL := fmt.Sprintf("%s/admin/reset-password/%s", s.config.ClientURL, resetPasswordToken)
	emailSubject := constant.ResetPasswordSubject
	emailContent := fmt.Sprintf("Please click the following url to reset password %s", resetPasswordURL)
	if err := s.emailService.SendEmail(emailSubject, user.Firstname, user.Email, emailContent); err != nil {
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	return &dto.ForgotPasswordResponse{
		IsSuccess: true,
	}, nil
}

func (s *serviceImpl) ResetPassword(request *dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, *dto.ResponseErr) {
	resetTokenCache, err := s.tokenService.FindResetPasswordToken(request.Token)
	if err != nil {
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	userDb := &model.User{}
	if err := s.userRepo.FindById(resetTokenCache.UserID, userDb); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.NotFoundError(constant.UserNotFoundErrorMessage)
		}
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	err = s.bcryptUtil.CompareHashedPassword(userDb.Password, request.Password)
	if err == nil {
		return nil, dto.BadRequestError(constant.IncorrectPasswordErrorMessage)
	}

	hashPassword, err := s.bcryptUtil.GenerateHashedPassword(request.Password)
	if err != nil {
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	userDb.Password = hashPassword
	if err := s.userRepo.Update(resetTokenCache.UserID, userDb); err != nil {
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	if err := s.tokenService.RemoveResetPasswordToken(request.Token); err != nil {
		return nil, dto.InternalServerError(constant.InternalServerErrorMessage)
	}

	return &dto.ResetPasswordResponse{
		IsSuccess: true,
	}, nil
}
