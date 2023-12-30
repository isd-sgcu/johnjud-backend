package auth

import (
	"errors"
	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	mock_router "github.com/isd-sgcu/johnjud-gateway/src/mocks/router"
	mock_auth "github.com/isd-sgcu/johnjud-gateway/src/mocks/service/auth"
	mock_user "github.com/isd-sgcu/johnjud-gateway/src/mocks/service/user"
	mock_validator "github.com/isd-sgcu/johnjud-gateway/src/mocks/validator"
	auth_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/auth/v1"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type AuthHandlerTest struct {
	suite.Suite
	signupRequest *dto.SignupRequest
	signInRequest *dto.SignIn
	bindErr       error
	validateErr   []*dto.BadReqErrResponse
}

func TestAuthHandler(t *testing.T) {
	suite.Run(t, new(AuthHandlerTest))
}

func (t *AuthHandlerTest) SetupTest() {
	signupRequest := &dto.SignupRequest{}
	signInRequest := &dto.SignIn{}
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

	authSvc := mock_auth.NewMockService(controller)
	userSvc := mock_user.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

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

	authSvc := mock_auth.NewMockService(controller)
	userSvc := mock_user.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

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

	authSvc := mock_auth.NewMockService(controller)
	userSvc := mock_user.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

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

	authSvc := mock_auth.NewMockService(controller)
	userSvc := mock_user.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Bind(t.signupRequest).Return(nil)
	validator.EXPECT().Validate(t.signupRequest).Return(nil)
	authSvc.EXPECT().Signup(t.signupRequest).Return(nil, signupError)
	context.EXPECT().JSON(http.StatusInternalServerError, signupError)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.Signup(context)
}

func (t *AuthHandlerTest) TestSignInSuccess() {
	signInResponse := &auth_proto.Credential{
		AccessToken:  faker.Word(),
		RefreshToken: faker.UUIDDigit(),
		ExpiresIn:    3600,
	}

	controller := gomock.NewController(t.T())

	authSvc := mock_auth.NewMockService(controller)
	userSvc := mock_user.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

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

	authSvc := mock_auth.NewMockService(controller)
	userSvc := mock_user.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

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

	authSvc := mock_auth.NewMockService(controller)
	userSvc := mock_user.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

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

	authSvc := mock_auth.NewMockService(controller)
	userSvc := mock_user.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Bind(t.signInRequest).Return(nil)
	validator.EXPECT().Validate(t.signInRequest).Return(nil)
	authSvc.EXPECT().SignIn(t.signInRequest).Return(nil, signInErrResponse)
	context.EXPECT().JSON(http.StatusInternalServerError, signInErrResponse)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.SignIn(context)
}
