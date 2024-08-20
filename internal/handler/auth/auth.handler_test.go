package auth

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	routerMock "github.com/isd-sgcu/johnjud-gateway/mocks/router"
	authMock "github.com/isd-sgcu/johnjud-gateway/mocks/service/auth"
	userMock "github.com/isd-sgcu/johnjud-gateway/mocks/service/user"
	validatorMock "github.com/isd-sgcu/johnjud-gateway/mocks/validator"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerTest struct {
	suite.Suite
	signupRequest         *dto.SignupRequest
	signInRequest         *dto.SignInRequest
	refreshTokenRequest   *dto.RefreshTokenRequest
	forgotPasswordRequest *dto.ForgotPasswordRequest
	resetPasswordRequest  *dto.ResetPasswordRequest
	bindErr               error
	validateErr           []*dto.BadReqErrResponse
}

func TestAuthHandler(t *testing.T) {
	suite.Run(t, new(AuthHandlerTest))
}

func (t *AuthHandlerTest) SetupTest() {
	signupRequest := &dto.SignupRequest{}
	signInRequest := &dto.SignInRequest{}
	refreshTokenRequest := &dto.RefreshTokenRequest{}
	forgotPasswordRequest := &dto.ForgotPasswordRequest{}
	resetPasswordRequest := &dto.ResetPasswordRequest{}
	bindErr := errors.New("Binding request failed")
	validateErr := []*dto.BadReqErrResponse{
		{
			Message:     "BadRequestError1",
			FailedField: "Field1",
			Value:       nil,
		},
		{
			Message:     "BadRequestError2",
			FailedField: "Field2",
			Value:       nil,
		},
	}

	t.signupRequest = signupRequest
	t.signInRequest = signInRequest
	t.refreshTokenRequest = refreshTokenRequest
	t.forgotPasswordRequest = forgotPasswordRequest
	t.resetPasswordRequest = resetPasswordRequest
	t.bindErr = bindErr
	t.validateErr = validateErr
}

func (t *AuthHandlerTest) TestSignupSuccess() {
	signupResponse := &dto.SignupResponse{
		Id:        faker.UUIDDigit(),
		Email:     faker.Email(),
		Firstname: faker.FirstName(),
		Lastname:  faker.LastName(),
	}
	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.signupRequest).Return(nil)
	validator.EXPECT().Validate(t.signupRequest).Return(nil)
	authSvc.EXPECT().Signup(t.signupRequest).Return(signupResponse, nil)
	context.EXPECT().JSON(http.StatusOK, signupResponse)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.Signup(context)
}

func (t *AuthHandlerTest) TestSignupBindFailed() {
	errResponse := dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.BindingRequestErrorMessage + t.bindErr.Error(),
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.signupRequest).Return(t.bindErr)
	context.EXPECT().JSON(http.StatusBadRequest, errResponse)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.Signup(context)
}

func (t *AuthHandlerTest) TestSignupValidateFailed() {
	errResponse := dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.InvalidRequestBodyMessage + "BadRequestError1, BadRequestError2",
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.signupRequest).Return(nil)
	validator.EXPECT().Validate(t.signupRequest).Return(t.validateErr)
	context.EXPECT().JSON(http.StatusBadRequest, errResponse)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.Signup(context)
}

func (t *AuthHandlerTest) TestSignupServiceError() {
	signupError := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    "Some Exception",
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.signupRequest).Return(nil)
	validator.EXPECT().Validate(t.signupRequest).Return(nil)
	authSvc.EXPECT().Signup(t.signupRequest).Return(nil, signupError)
	context.EXPECT().JSON(http.StatusInternalServerError, signupError)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.Signup(context)
}

func (t *AuthHandlerTest) TestSignInSuccess() {
	signInResponse := &dto.Credential{
		AccessToken:  faker.Word(),
		RefreshToken: faker.UUIDDigit(),
		ExpiresIn:    3600,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.signInRequest).Return(nil)
	validator.EXPECT().Validate(t.signInRequest).Return(nil)
	authSvc.EXPECT().SignIn(t.signInRequest).Return(signInResponse, nil)
	context.EXPECT().JSON(http.StatusOK, signInResponse)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.SignIn(context)
}

func (t *AuthHandlerTest) TestSignInBindFailed() {
	errResponse := dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.BindingRequestErrorMessage + t.bindErr.Error(),
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.signInRequest).Return(t.bindErr)
	context.EXPECT().JSON(http.StatusBadRequest, errResponse)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.SignIn(context)
}

func (t *AuthHandlerTest) TestSignInValidateFailed() {
	errResponse := dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.InvalidRequestBodyMessage + "BadRequestError1, BadRequestError2",
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.signInRequest).Return(nil)
	validator.EXPECT().Validate(t.signInRequest).Return(t.validateErr)
	context.EXPECT().JSON(http.StatusBadRequest, errResponse)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.SignIn(context)
}

func (t *AuthHandlerTest) TestSignInServiceError() {
	signInErrResponse := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.signInRequest).Return(nil)
	validator.EXPECT().Validate(t.signInRequest).Return(nil)
	authSvc.EXPECT().SignIn(t.signInRequest).Return(nil, signInErrResponse)
	context.EXPECT().JSON(http.StatusInternalServerError, signInErrResponse)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.SignIn(context)
}

func (t *AuthHandlerTest) TestSignOutSuccess() {
	token := faker.Word()
	signOutResponse := &dto.SignOutResponse{
		IsSuccess: true,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Token().Return(token)
	authSvc.EXPECT().SignOut(token).Return(signOutResponse, nil)
	context.EXPECT().JSON(http.StatusOK, signOutResponse)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.SignOut(context)
}

func (t *AuthHandlerTest) TestSignOutServiceError() {
	token := faker.Word()
	errResponse := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	handler := NewHandler(authSvc, userSvc, validator)

	context.EXPECT().Token().Return(token)
	authSvc.EXPECT().SignOut(token).Return(nil, errResponse)
	context.EXPECT().JSON(http.StatusInternalServerError, errResponse)

	handler.SignOut(context)
}

func (t *AuthHandlerTest) TestRefreshTokenSuccess() {
	refreshTokenResponse := &dto.Credential{
		AccessToken:  faker.Word(),
		RefreshToken: faker.UUIDDigit(),
		ExpiresIn:    3600,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	handler := NewHandler(authSvc, userSvc, validator)

	context.EXPECT().Bind(t.refreshTokenRequest).Return(nil)
	validator.EXPECT().Validate(t.refreshTokenRequest).Return(nil)
	authSvc.EXPECT().RefreshToken(t.refreshTokenRequest).Return(refreshTokenResponse, nil)
	context.EXPECT().JSON(http.StatusOK, refreshTokenResponse)

	handler.RefreshToken(context)
}

func (t *AuthHandlerTest) TestRefreshTokenBindFailed() {
	errResponse := dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.BindingRequestErrorMessage + t.bindErr.Error(),
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	handler := NewHandler(authSvc, userSvc, validator)

	context.EXPECT().Bind(t.refreshTokenRequest).Return(t.bindErr)
	context.EXPECT().JSON(http.StatusBadRequest, errResponse)

	handler.RefreshToken(context)
}

func (t *AuthHandlerTest) TestRefreshTokenValidateFailed() {
	errResponse := dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.InvalidRequestBodyMessage + "BadRequestError1, BadRequestError2",
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	handler := NewHandler(authSvc, userSvc, validator)

	context.EXPECT().Bind(t.refreshTokenRequest).Return(nil)
	validator.EXPECT().Validate(t.refreshTokenRequest).Return(t.validateErr)
	context.EXPECT().JSON(http.StatusBadRequest, errResponse)

	handler.RefreshToken(context)
}

func (t *AuthHandlerTest) TestRefreshTokenServiceError() {
	refreshTokenErr := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	handler := NewHandler(authSvc, userSvc, validator)

	context.EXPECT().Bind(t.refreshTokenRequest).Return(nil)
	validator.EXPECT().Validate(t.refreshTokenRequest).Return(nil)
	authSvc.EXPECT().RefreshToken(t.refreshTokenRequest).Return(nil, refreshTokenErr)
	context.EXPECT().JSON(http.StatusInternalServerError, refreshTokenErr)

	handler.RefreshToken(context)
}

func (t *AuthHandlerTest) TestForgotPasswordSuccess() {
	forgotPasswordResponse := &dto.ForgotPasswordResponse{
		IsSuccess: true,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	handler := NewHandler(authSvc, userSvc, validator)

	context.EXPECT().Bind(t.forgotPasswordRequest).Return(nil)
	validator.EXPECT().Validate(t.forgotPasswordRequest).Return(nil)
	authSvc.EXPECT().ForgotPassword(t.forgotPasswordRequest).Return(forgotPasswordResponse, nil)
	context.EXPECT().JSON(http.StatusOK, forgotPasswordResponse)

	handler.ForgotPassword(context)
}

func (t *AuthHandlerTest) TestForgotPasswordBindFailed() {
	errResponse := dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.BindingRequestErrorMessage + t.bindErr.Error(),
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	handler := NewHandler(authSvc, userSvc, validator)

	context.EXPECT().Bind(t.forgotPasswordRequest).Return(t.bindErr)
	context.EXPECT().JSON(http.StatusBadRequest, errResponse)
	handler.ForgotPassword(context)
}

func (t *AuthHandlerTest) TestForgotPasswordValidateFailed() {
	errResponse := dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.InvalidRequestBodyMessage + "BadRequestError1, BadRequestError2",
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	handler := NewHandler(authSvc, userSvc, validator)

	context.EXPECT().Bind(t.forgotPasswordRequest).Return(nil)
	validator.EXPECT().Validate(t.forgotPasswordRequest).Return(t.validateErr)
	context.EXPECT().JSON(http.StatusBadRequest, errResponse)

	handler.ForgotPassword(context)
}

func (t *AuthHandlerTest) TestForgotPasswordServiceError() {
	forgotPasswordErr := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	handler := NewHandler(authSvc, userSvc, validator)

	context.EXPECT().Bind(t.forgotPasswordRequest).Return(nil)
	validator.EXPECT().Validate(t.forgotPasswordRequest).Return(nil)
	authSvc.EXPECT().ForgotPassword(t.forgotPasswordRequest).Return(nil, forgotPasswordErr)
	context.EXPECT().JSON(http.StatusInternalServerError, forgotPasswordErr)

	handler.ForgotPassword(context)
}

func (t *AuthHandlerTest) TestResetPasswordSuccess() {
	resetPasswordResponse := &dto.ResetPasswordResponse{
		IsSuccess: true,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	handler := NewHandler(authSvc, userSvc, validator)

	context.EXPECT().Bind(t.resetPasswordRequest).Return(nil)
	validator.EXPECT().Validate(t.resetPasswordRequest).Return(nil)
	authSvc.EXPECT().ResetPassword(t.resetPasswordRequest).Return(resetPasswordResponse, nil)
	context.EXPECT().JSON(http.StatusOK, resetPasswordResponse)

	handler.ResetPassword(context)
}

func (t *AuthHandlerTest) TestResetPasswordBindFailed() {
	errResponse := dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.BindingRequestErrorMessage + t.bindErr.Error(),
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	handler := NewHandler(authSvc, userSvc, validator)

	context.EXPECT().Bind(t.resetPasswordRequest).Return(t.bindErr)
	context.EXPECT().JSON(http.StatusBadRequest, errResponse)

	handler.ResetPassword(context)
}

func (t *AuthHandlerTest) TestResetPasswordValidateFailed() {
	errResponse := dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.InvalidRequestBodyMessage + "BadRequestError1, BadRequestError2",
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	handler := NewHandler(authSvc, userSvc, validator)

	context.EXPECT().Bind(t.resetPasswordRequest).Return(nil)
	validator.EXPECT().Validate(t.resetPasswordRequest).Return(t.validateErr)
	context.EXPECT().JSON(http.StatusBadRequest, errResponse)

	handler.ResetPassword(context)
}

func (t *AuthHandlerTest) TestResetPasswordServiceError() {
	resetPasswordErr := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := authMock.NewMockService(controller)
	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	handler := NewHandler(authSvc, userSvc, validator)

	context.EXPECT().Bind(t.resetPasswordRequest).Return(nil)
	validator.EXPECT().Validate(t.resetPasswordRequest).Return(nil)
	authSvc.EXPECT().ResetPassword(t.resetPasswordRequest).Return(nil, resetPasswordErr)
	context.EXPECT().JSON(http.StatusInternalServerError, resetPasswordErr)

	handler.ResetPassword(context)
}
