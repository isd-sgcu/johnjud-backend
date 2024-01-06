package like

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	utils "github.com/isd-sgcu/johnjud-gateway/src/app/utils/like"
	likeMock "github.com/isd-sgcu/johnjud-gateway/src/mocks/client/like"
	likeProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/like/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LikeServiceTest struct {
	suite.Suite
	Likes              []*likeProto.Like
	Like               *likeProto.Like
	LikeResponse       *dto.LikeResponse
	CreateLikeProtoReq *likeProto.CreateLikeRequest
	CreateLikeDtoReq   *dto.CreateLikeRequest
	DeleteLikeProtoReq *likeProto.DeleteLikeRequest
	DeleteLikeDtoReq   *dto.DeleteLikeRequest
	NotFoundErr        *dto.ResponseErr

	UnavailableServiceErr *dto.ResponseErr
	InvalidArgumentErr    *dto.ResponseErr
	InternalErr           *dto.ResponseErr
}

func TestLikeService(t *testing.T) {
	suite.Run(t, new(LikeServiceTest))
}

func (t *LikeServiceTest) SetupTest() {
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

	t.CreateLikeProtoReq = &likeProto.CreateLikeRequest{
		Like: &likeProto.Like{
			UserId: t.Like.UserId,
			PetId:  t.Like.PetId,
		},
	}

	t.CreateLikeDtoReq = &dto.CreateLikeRequest{
		UserID: t.Like.UserId,
		PetID:  t.Like.PetId,
	}

	t.DeleteLikeProtoReq = &likeProto.DeleteLikeRequest{
		Id: t.Like.Id,
	}

	t.UnavailableServiceErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    constant.UnavailableServiceMessage,
		Data:       nil,
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    constant.UserNotFoundMessage,
		Data:       nil,
	}

	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	t.InvalidArgumentErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.InvalidArgumentMessage,
		Data:       nil,
	}
}

func (t *LikeServiceTest) TestFindByUserIdSuccess() {
	protoReq := &likeProto.FindLikeByUserIdRequest{
		UserId: t.Like.UserId,
	}
	protoResp := &likeProto.FindLikeByUserIdResponse{
		Likes: t.Likes,
	}

	expected := utils.ProtoToDtoList(t.Likes)

	client := likeMock.LikeClientMock{}
	client.On("FindByUserId", protoReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.FindByUserId(t.Like.UserId)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *LikeServiceTest) TestFindByUserIdNotFoundError() {
	protoReq := &likeProto.FindLikeByUserIdRequest{
		UserId: t.Like.UserId,
	}

	clientErr := status.Error(codes.NotFound, constant.UserNotFoundMessage)

	expected := t.NotFoundErr

	client := likeMock.LikeClientMock{}
	client.On("FindByUserId", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.FindByUserId(t.Like.UserId)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *LikeServiceTest) TestFindByUserIdUnavailableServiceError() {
	protoReq := &likeProto.FindLikeByUserIdRequest{
		UserId: t.Like.UserId,
	}

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := likeMock.LikeClientMock{}
	client.On("FindByUserId", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.FindByUserId(t.Like.UserId)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *LikeServiceTest) TestFindByUserIdInternalError() {
	protoReq := &likeProto.FindLikeByUserIdRequest{
		UserId: t.Like.UserId,
	}

	clientErr := status.Error(codes.Internal, constant.InternalErrorMessage)

	expected := t.InternalErr

	client := likeMock.LikeClientMock{}
	client.On("FindByUserId", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.FindByUserId(t.Like.UserId)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *LikeServiceTest) TestCreateSuccess() {
	protoReq := t.CreateLikeProtoReq
	protoResp := &likeProto.CreateLikeResponse{
		Like: t.Like,
	}

	expected := utils.ProtoToDto(t.Like)

	client := &likeMock.LikeClientMock{}
	client.On("Create", protoReq).Return(protoResp, nil)

	svc := NewService(client)
	actual, err := svc.Create(t.CreateLikeDtoReq)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *LikeServiceTest) TestCreateInvalidArgumentError() {
	protoReq := t.CreateLikeProtoReq

	expected := t.InvalidArgumentErr

	clientErr := status.Error(codes.InvalidArgument, constant.InvalidArgumentMessage)

	client := &likeMock.LikeClientMock{}
	client.On("Create", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Create(t.CreateLikeDtoReq)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *LikeServiceTest) TestCreateInternalError() {
	protoReq := t.CreateLikeProtoReq

	expected := t.InvalidArgumentErr

	clientErr := status.Error(codes.InvalidArgument, constant.InvalidArgumentMessage)

	client := &likeMock.LikeClientMock{}
	client.On("Create", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Create(t.CreateLikeDtoReq)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *LikeServiceTest) TestCreateUnavailableServiceError() {
	protoReq := t.CreateLikeProtoReq

	expected := t.UnavailableServiceErr

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	client := &likeMock.LikeClientMock{}
	client.On("Create", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Create(t.CreateLikeDtoReq)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *LikeServiceTest) TestDeleteSuccess() {
	protoReq := &likeProto.DeleteLikeRequest{
		Id: t.Like.Id,
	}
	protoResp := &likeProto.DeleteLikeResponse{
		Success: true,
	}

	expected := &dto.DeleteLikeResponse{Success: true}

	client := &likeMock.LikeClientMock{}
	client.On("Delete", protoReq).Return(protoResp, nil)

	svc := NewService(client)
	actual, err := svc.Delete(t.Like.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *LikeServiceTest) TestDeleteNotFoundError() {
	protoReq := &likeProto.DeleteLikeRequest{
		Id: t.Like.Id,
	}

	clientErr := status.Error(codes.NotFound, constant.UserNotFoundMessage)

	expected := t.NotFoundErr

	client := &likeMock.LikeClientMock{}
	client.On("Delete", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Delete(t.Like.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *LikeServiceTest) TestDeleteUnavailableServiceError() {
	protoReq := &likeProto.DeleteLikeRequest{
		Id: t.Like.Id,
	}

	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)

	expected := t.UnavailableServiceErr

	client := &likeMock.LikeClientMock{}
	client.On("Delete", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Delete(t.Like.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *LikeServiceTest) TestDeleteInternalError() {
	protoReq := &likeProto.DeleteLikeRequest{
		Id: t.Like.Id,
	}

	clientErr := status.Error(codes.Internal, constant.InternalErrorMessage)

	expected := t.InternalErr

	client := &likeMock.LikeClientMock{}
	client.On("Delete", protoReq).Return(nil, clientErr)

	svc := NewService(client)
	actual, err := svc.Delete(t.Like.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}
