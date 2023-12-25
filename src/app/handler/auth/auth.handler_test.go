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
}

func TestAuthHandler(t *testing.T) {
	suite.Run(t, new(AuthHandlerTest))
}

func (t *AuthHandlerTest) SetupTest() {

}

func (t *AuthHandlerTest) TestSignupSuccess() {
	signupRequest := &dto.SignupRequest{}
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

	context.EXPECT().Bind(signupRequest).Return(nil)
	validator.EXPECT().Validate(signupRequest).Return(nil)
	authSvc.EXPECT().Signup(signupRequest).Return(signupResponse, nil)
	context.EXPECT().JSON(http.StatusOK, signupResponse)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.Signup(context)
}

func (t *AuthHandlerTest) TestSignupBindFailed() {
	signupRequest := &dto.SignupRequest{}
	bindReqErr := errors.New("Binding request failed")
	errResponse := dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.BindingRequestErrorMessage + bindReqErr.Error(),
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	authSvc := mock_auth.NewMockService(controller)
	userSvc := mock_user.NewMockService(controller)
	validator := mock_validator.NewMockIDtoValidator(controller)
	context := mock_router.NewMockIContext(controller)

	context.EXPECT().Bind(signupRequest).Return(bindReqErr)
	context.EXPECT().JSON(http.StatusBadRequest, errResponse)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.Signup(context)
}

func (t *AuthHandlerTest) TestSignupValidateFailed() {
	signupRequest := &dto.SignupRequest{}
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

	context.EXPECT().Bind(signupRequest).Return(nil)
	validator.EXPECT().Validate(signupRequest).Return(validateErr)
	context.EXPECT().JSON(http.StatusBadRequest, errResponse)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.Signup(context)
}

func (t *AuthHandlerTest) TestSignupServiceError() {
	signupRequest := &dto.SignupRequest{}
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

	context.EXPECT().Bind(signupRequest).Return(nil)
	validator.EXPECT().Validate(signupRequest).Return(nil)
	authSvc.EXPECT().Signup(signupRequest).Return(nil, signupError)
	context.EXPECT().JSON(http.StatusInternalServerError, signupError)

	handler := NewHandler(authSvc, userSvc, validator)

	handler.Signup(context)
}
