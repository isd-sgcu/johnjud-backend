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
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type AuthHandlerTest struct {
	suite.Suite
	signupRequest *dto.SignupRequest
	bindErr       error
	validateErr   []*dto.BadReqErrResponse
}

func TestAuthHandler(t *testing.T) {
	suite.Run(t, new(AuthHandlerTest))
}

func (t *AuthHandlerTest) SetupTest() {
	signupRequest := &dto.SignupRequest{}
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
