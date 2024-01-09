package user

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	routerMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/router"
	userMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/service/user"
	validatorMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/validator"

	errConst "github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"

	"github.com/stretchr/testify/suite"
)

type UserHandlerTest struct {
	suite.Suite
	User              *proto.User
	UserDto           *dto.User
	UpdateUserRequest *dto.UpdateUserRequest
	BindErr           *dto.ResponseErr
	NotFoundErr       *dto.ResponseErr
	InvalidIDErr      *dto.ResponseErr
	ServiceDownErr    *dto.ResponseErr
	DuplicateEmailErr *dto.ResponseErr
	InternalErr       *dto.ResponseErr
}

func TestUserHandler(t *testing.T) {
	suite.Run(t, new(UserHandlerTest))
}

func (t *UserHandlerTest) SetupTest() {
	t.User = &proto.User{
		Id:        faker.UUIDDigit(),
		Email:     faker.Email(),
		Password:  faker.Password(),
		Firstname: faker.FirstName(),
		Lastname:  faker.LastName(),
		Role:      "user",
	}
	t.UserDto = &dto.User{
		Id:        t.User.Id,
		Firstname: t.User.Firstname,
		Lastname:  t.User.Lastname,
		Email:     t.User.Email,
	}

	t.UpdateUserRequest = &dto.UpdateUserRequest{}

	t.ServiceDownErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service is down",
		Data:       nil,
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    errConst.UserNotFoundMessage,
		Data:       nil,
	}

	t.InvalidIDErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    errConst.InvalidIDMessage,
		Data:       nil,
	}

	t.BindErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    errConst.InvalidIDMessage,
	}

	t.DuplicateEmailErr = &dto.ResponseErr{
		StatusCode: http.StatusConflict,
		Message:    errConst.DuplicateEmailMessage,
		Data:       nil,
	}

	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    errConst.InternalErrorMessage,
		Data:       nil,
	}
}

func (t *UserHandlerTest) TestFindOneSuccess() {
	svcResp := t.UserDto
	expectedResp := t.UserDto

	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().ID().Return(t.User.Id, nil)
	userSvc.EXPECT().FindOne(t.User.Id).Return(svcResp, nil)
	context.EXPECT().JSON(http.StatusOK, expectedResp)

	handler := NewHandler(userSvc, validator)
	handler.FindOne(context)
}

func (t *UserHandlerTest) TestFindOneNotFoundErr() {
	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().ID().Return(t.User.Id, nil)
	userSvc.EXPECT().FindOne(t.User.Id).Return(nil, t.NotFoundErr)
	context.EXPECT().JSON(http.StatusNotFound, t.NotFoundErr)

	handler := NewHandler(userSvc, validator)
	handler.FindOne(context)
}

func (t *UserHandlerTest) TestFindOneInvalidIDErr() {
	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().ID().Return("", errors.New("Invalid ID"))
	context.EXPECT().JSON(http.StatusBadRequest, t.InvalidIDErr)

	handler := NewHandler(userSvc, validator)
	handler.FindOne(context)
}

func (t *UserHandlerTest) TestFindOneInternalErr() {
	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().ID().Return(t.User.Id, nil)
	userSvc.EXPECT().FindOne(t.User.Id).Return(nil, t.InternalErr)
	context.EXPECT().JSON(http.StatusInternalServerError, t.InternalErr)

	handler := NewHandler(userSvc, validator)
	handler.FindOne(context)
}

func (t *UserHandlerTest) TestFindOneGrpcErr() {
	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().ID().Return(t.User.Id, nil)
	userSvc.EXPECT().FindOne(t.User.Id).Return(nil, t.ServiceDownErr)
	context.EXPECT().JSON(http.StatusServiceUnavailable, t.ServiceDownErr)

	handler := NewHandler(userSvc, validator)
	handler.FindOne(context)
}

func (t *UserHandlerTest) TestUpdateSuccess() {
	svcResp := t.UserDto
	expectedResp := t.UserDto

	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().UserID().Return(t.User.Id)
	context.EXPECT().Bind(t.UpdateUserRequest).Return(nil)
	validator.EXPECT().Validate(t.UpdateUserRequest).Return(nil)
	userSvc.EXPECT().Update(t.User.Id, t.UpdateUserRequest).Return(svcResp, nil)
	context.EXPECT().JSON(http.StatusOK, expectedResp)

	handler := NewHandler(userSvc, validator)
	handler.Update(context)
}

func (t *UserHandlerTest) TestUpdateBindErr() {
	bindErr := errors.New("Bind err")
	expectedResp := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.BindingRequestErrorMessage + bindErr.Error(),
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().UserID().Return(t.User.Id)
	context.EXPECT().Bind(t.UpdateUserRequest).Return(bindErr)
	context.EXPECT().JSON(http.StatusBadRequest, expectedResp)

	handler := NewHandler(userSvc, validator)
	handler.Update(context)
}

func (t *UserHandlerTest) TestUpdateValidateErr() {
	errorMessage := []string{"First name is required", "Last name is required"}
	validateErr := []*dto.BadReqErrResponse{
		{
			Message:     errorMessage[0],
			FailedField: "firstname",
		},
		{
			Message:     errorMessage[1],
			FailedField: "lastname",
		},
	}
	expectedResp := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.InvalidRequestBodyMessage + strings.Join(errorMessage, ", "),
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().UserID().Return(t.User.Id)
	context.EXPECT().Bind(t.UpdateUserRequest).Return(nil)
	validator.EXPECT().Validate(t.UpdateUserRequest).Return(validateErr)
	context.EXPECT().JSON(http.StatusBadRequest, expectedResp)

	handler := NewHandler(userSvc, validator)
	handler.Update(context)
}

func (t *UserHandlerTest) TestUpdateDuplicateEmailErr() {
	expectedResp := &dto.ResponseErr{
		StatusCode: http.StatusConflict,
		Message:    constant.DuplicateEmailMessage,
		Data:       nil,
	}

	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().UserID().Return(t.User.Id)
	context.EXPECT().Bind(t.UpdateUserRequest).Return(nil)
	validator.EXPECT().Validate(t.UpdateUserRequest).Return(nil)
	userSvc.EXPECT().Update(t.User.Id, t.UpdateUserRequest).Return(nil, t.DuplicateEmailErr)
	context.EXPECT().JSON(http.StatusConflict, expectedResp)

	handler := NewHandler(userSvc, validator)
	handler.Update(context)
}

func (t *UserHandlerTest) TestUpdateInternalErr() {
	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().UserID().Return(t.User.Id)
	context.EXPECT().Bind(t.UpdateUserRequest).Return(nil)
	validator.EXPECT().Validate(t.UpdateUserRequest).Return(nil)
	userSvc.EXPECT().Update(t.User.Id, t.UpdateUserRequest).Return(nil, t.InternalErr)
	context.EXPECT().JSON(http.StatusInternalServerError, t.InternalErr)

	handler := NewHandler(userSvc, validator)
	handler.Update(context)
}

func (t *UserHandlerTest) TestUpdateGrpcErr() {
	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().UserID().Return(t.User.Id)
	context.EXPECT().Bind(t.UpdateUserRequest).Return(nil)
	validator.EXPECT().Validate(t.UpdateUserRequest).Return(nil)
	userSvc.EXPECT().Update(t.User.Id, t.UpdateUserRequest).Return(nil, t.ServiceDownErr)
	context.EXPECT().JSON(http.StatusServiceUnavailable, t.ServiceDownErr)

	handler := NewHandler(userSvc, validator)
	handler.Update(context)
}

func (t *UserHandlerTest) TestDeleteSuccess() {
	deleteResp := &dto.DeleteUserResponse{
		Success: true,
	}
	expectedResp := deleteResp

	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().ID().Return(t.User.Id, nil)
	userSvc.EXPECT().Delete(t.User.Id).Return(deleteResp, nil)
	context.EXPECT().JSON(http.StatusOK, expectedResp)

	handler := NewHandler(userSvc, validator)
	handler.Delete(context)
}

func (t *UserHandlerTest) TestDeleteInvalidIDErr() {
	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().ID().Return("", errors.New("Invalid ID"))
	context.EXPECT().JSON(http.StatusBadRequest, t.InvalidIDErr)

	handler := NewHandler(userSvc, validator)
	handler.Delete(context)
}

func (t *UserHandlerTest) TestDeleteInternalErr() {
	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().ID().Return(t.User.Id, nil)
	userSvc.EXPECT().Delete(t.User.Id).Return(nil, t.InternalErr)
	context.EXPECT().JSON(http.StatusInternalServerError, t.InternalErr)

	handler := NewHandler(userSvc, validator)
	handler.Delete(context)
}

func (t *UserHandlerTest) TestDeleteGrpcErr() {
	controller := gomock.NewController(t.T())

	userSvc := userMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().ID().Return(t.User.Id, nil)
	userSvc.EXPECT().Delete(t.User.Id).Return(nil, t.ServiceDownErr)
	context.EXPECT().JSON(http.StatusServiceUnavailable, t.ServiceDownErr)

	handler := NewHandler(userSvc, validator)
	handler.Delete(context)
}
