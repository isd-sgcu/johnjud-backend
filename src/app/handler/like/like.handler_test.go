package auth

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/golang/mock/gomock"
	errConst "github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	utils "github.com/isd-sgcu/johnjud-gateway/src/app/utils/like"
	likeConst "github.com/isd-sgcu/johnjud-gateway/src/constant/like"
	routerMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/router"
	likeMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/service/like"
	validatorMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/validator"
	likeProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/like/v1"
	"github.com/stretchr/testify/suite"
)

type LikeHandlerTest struct {
	suite.Suite
	Likes             []*likeProto.Like
	Like              *likeProto.Like
	LikeResponse      *dto.LikeResponse
	CreateLikeRequest *dto.CreateLikeRequest
	DeleteLikeRequest *dto.DeleteLikeRequest
	NotFoundErr       *dto.ResponseErr

	UnavailableServiceErr *dto.ResponseErr
	InvalidArgumentErr    *dto.ResponseErr
	BindErr               *dto.ResponseErr
	InternalErr           *dto.ResponseErr
}

func TestLikeHandler(t *testing.T) {
	suite.Run(t, new(LikeHandlerTest))
}

func (t *LikeHandlerTest) SetupTest() {
	var likes []*likeProto.Like
	for i := 0; i <= 3; i++ {
		like := &likeProto.Like{
			Id:     faker.UUIDDigit(),
			UserId: faker.UUIDDigit(),
			PetId:  faker.UUIDDigit(),
		}
		likes = append(likes, like)
	}

	t.Likes = likes
	t.Like = likes[0]

	t.CreateLikeRequest = &dto.CreateLikeRequest{}
	t.DeleteLikeRequest = &dto.DeleteLikeRequest{}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    errConst.UserNotFoundMessage,
		Data:       nil,
	}

	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    errConst.InternalErrorMessage,
		Data:       nil,
	}

	t.UnavailableServiceErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    errConst.UnavailableServiceMessage,
		Data:       nil,
	}
}

func (t *LikeHandlerTest) TestFindLikesSuccess() {
	findLikeResponse := utils.ProtoToDtoList(t.Likes)
	expectedResponse := dto.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    likeConst.FindLikeSuccessMessage,
		Data:       findLikeResponse,
	}

	controller := gomock.NewController(t.T())

	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Like.UserId, nil)
	likeSvc.EXPECT().FindByUserId(t.Like.UserId).Return(utils.ProtoToDtoList(t.Likes), nil)
	context.EXPECT().JSON(http.StatusOK, expectedResponse)

	handler := NewHandler(likeSvc, validator)
	handler.FindByUserId(context)
}

func (t *LikeHandlerTest) TestFindLikeNotFoundError() {
	findLikeErrorResponse := t.NotFoundErr

	controller := gomock.NewController(t.T())

	petSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Like.UserId, nil)
	petSvc.EXPECT().FindByUserId(t.Like.UserId).Return(nil, findLikeErrorResponse)
	context.EXPECT().JSON(http.StatusNotFound, findLikeErrorResponse)

	handler := NewHandler(petSvc, validator)
	handler.FindByUserId(context)
}

func (t *LikeHandlerTest) TestFindLikeServiceUnavailableError() {
	findLikeErrorResponse := t.UnavailableServiceErr

	controller := gomock.NewController(t.T())

	petSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Like.UserId, nil)
	petSvc.EXPECT().FindByUserId(t.Like.UserId).Return(nil, findLikeErrorResponse)
	context.EXPECT().JSON(http.StatusServiceUnavailable, findLikeErrorResponse)

	handler := NewHandler(petSvc, validator)
	handler.FindByUserId(context)
}

func (t *LikeHandlerTest) TestFindLikeInternalError() {
	findLikeErrorResponse := t.InternalErr

	controller := gomock.NewController(t.T())

	petSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Like.UserId, nil)
	petSvc.EXPECT().FindByUserId(t.Like.UserId).Return(nil, findLikeErrorResponse)
	context.EXPECT().JSON(http.StatusInternalServerError, findLikeErrorResponse)

	handler := NewHandler(petSvc, validator)
	handler.FindByUserId(context)
}

func (t *LikeHandlerTest) TestCreateSuccess() {
	createLikeResponse := utils.ProtoToDto(t.Like)
	expectedResponse := dto.ResponseSuccess{
		StatusCode: http.StatusCreated,
		Message:    likeConst.CreateLikeSuccessMessage,
		Data:       createLikeResponse,
	}

	controller := gomock.NewController(t.T())

	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.CreateLikeRequest).Return(nil)
	validator.EXPECT().Validate(t.CreateLikeRequest).Return(nil)
	likeSvc.EXPECT().Create(t.CreateLikeRequest).Return(createLikeResponse, nil)
	context.EXPECT().JSON(http.StatusCreated, expectedResponse)

	handler := NewHandler(likeSvc, validator)
	handler.Create(context)
}

func (t *LikeHandlerTest) TestCreateUnavailableServiceError() {
	createLikeErrorResponse := t.UnavailableServiceErr

	controller := gomock.NewController(t.T())

	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.CreateLikeRequest).Return(nil)
	validator.EXPECT().Validate(t.CreateLikeRequest).Return(nil)
	likeSvc.EXPECT().Create(t.CreateLikeRequest).Return(nil, createLikeErrorResponse)
	context.EXPECT().JSON(http.StatusServiceUnavailable, createLikeErrorResponse)

	handler := NewHandler(likeSvc, validator)
	handler.Create(context)
}

func (t *LikeHandlerTest) TestCreateInternalError() {
	createLikeErrorResponse := t.InternalErr

	controller := gomock.NewController(t.T())

	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Bind(t.CreateLikeRequest).Return(nil)
	validator.EXPECT().Validate(t.CreateLikeRequest).Return(nil)
	likeSvc.EXPECT().Create(t.CreateLikeRequest).Return(nil, createLikeErrorResponse)
	context.EXPECT().JSON(http.StatusInternalServerError, createLikeErrorResponse)

	handler := NewHandler(likeSvc, validator)
	handler.Create(context)
}

func (t *LikeHandlerTest) TestDeleteSuccess() {
	deleteResponse := &dto.DeleteLikeResponse{
		Success: true,
	}
	expectedResponse := dto.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    likeConst.DelteLikeSuccessMessage,
		Data:       deleteResponse,
	}

	controller := gomock.NewController(t.T())

	likeSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Like.UserId, nil)
	likeSvc.EXPECT().Delete(t.Like.UserId).Return(deleteResponse, nil)
	context.EXPECT().JSON(http.StatusOK, expectedResponse)

	handler := NewHandler(likeSvc, validator)
	handler.Delete(context)
}

func (t *LikeHandlerTest) TestDeleteNotFoundError() {
	deleteErrorResponse := t.NotFoundErr

	controller := gomock.NewController(t.T())

	petSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Like.UserId, nil)
	petSvc.EXPECT().FindByUserId(t.Like.UserId).Return(nil, deleteErrorResponse)
	context.EXPECT().JSON(http.StatusNotFound, deleteErrorResponse)

	handler := NewHandler(petSvc, validator)
	handler.FindByUserId(context)
}

func (t *LikeHandlerTest) TestDeleteServiceUnavailableError() {
	deleteErrorResponse := t.UnavailableServiceErr

	controller := gomock.NewController(t.T())

	petSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Like.UserId, nil)
	petSvc.EXPECT().FindByUserId(t.Like.UserId).Return(nil, deleteErrorResponse)
	context.EXPECT().JSON(http.StatusServiceUnavailable, deleteErrorResponse)

	handler := NewHandler(petSvc, validator)
	handler.FindByUserId(context)
}

func (t *LikeHandlerTest) TestDeleteInternalError() {
	deleteErrorResponse := t.InternalErr

	controller := gomock.NewController(t.T())

	petSvc := likeMock.NewMockService(controller)
	validator := validatorMock.NewMockIDtoValidator(controller)
	context := routerMock.NewMockIContext(controller)

	context.EXPECT().Param("id").Return(t.Like.UserId, nil)
	petSvc.EXPECT().FindByUserId(t.Like.UserId).Return(nil, deleteErrorResponse)
	context.EXPECT().JSON(http.StatusInternalServerError, deleteErrorResponse)

	handler := NewHandler(petSvc, validator)
	handler.FindByUserId(context)
}
