package user

import (
	"net/http"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/mocks/client/user"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceTest struct {
	suite.Suite
	User                  *proto.User
	FindOneUserReq        *proto.FindOneUserRequest
	UpdateUserReq         *proto.UpdateUserRequest
	UpdateUserDto         *dto.UpdateUserRequest
	DeleteUserReq         *proto.DeleteUserRequest
	NotFoundErr           *dto.ResponseErr
	UnavailableServiceErr *dto.ResponseErr
	ConflictErr           *dto.ResponseErr
	InternalErr           *dto.ResponseErr
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTest))
}

func (t *UserServiceTest) SetupTest() {
	t.User = &proto.User{
		Id:        faker.UUIDDigit(),
		Email:     faker.Email(),
		Password:  faker.Password(),
		Firstname: faker.FirstName(),
		Lastname:  faker.LastName(),
		Role:      "user",
	}

	t.FindOneUserReq = &proto.FindOneUserRequest{
		Id: t.User.Id,
	}

	t.UpdateUserDto = &dto.UpdateUserRequest{
		Email:     faker.Email(),
		Password:  faker.Password(),
		Firstname: faker.FirstName(),
		Lastname:  faker.LastName(),
	}

	t.UpdateUserReq = &proto.UpdateUserRequest{
		Id:        t.User.Id,
		Email:     t.UpdateUserDto.Email,
		Password:  t.UpdateUserDto.Password,
		Firstname: t.UpdateUserDto.Firstname,
		Lastname:  t.UpdateUserDto.Lastname,
	}

	t.DeleteUserReq = &proto.DeleteUserRequest{
		Id: t.User.Id,
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

	t.ConflictErr = &dto.ResponseErr{
		StatusCode: http.StatusConflict,
		Message:    constant.DuplicateEmailMessage,
		Data:       nil,
	}

	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}
}

func (t *UserServiceTest) TestFindOneSuccess() {
	protoResp := &proto.FindOneUserResponse{
		User: &proto.User{
			Id:        t.User.Id,
			Email:     t.User.Email,
			Firstname: t.User.Firstname,
			Lastname:  t.User.Lastname,
			Role:      t.User.Role,
		},
	}

	expected := &dto.FindOneUserResponse{
		Id:        t.User.Id,
		Email:     t.User.Email,
		Firstname: t.User.Firstname,
		Lastname:  t.User.Lastname,
	}

	client := user.UserClientMock{}
	client.On("FindOne", t.FindOneUserReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.FindOne(t.User.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *UserServiceTest) TestFindOneNotFoundError() {
	expected := t.NotFoundErr

	client := user.UserClientMock{}
	clienErr := status.Error(codes.NotFound, constant.UserNotFoundMessage)
	client.On("FindOne", t.FindOneUserReq).Return(nil, clienErr)

	svc := NewService(&client)
	actual, err := svc.FindOne(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *UserServiceTest) TestFindOneUnavailableServiceError() {
	expected := t.UnavailableServiceErr

	client := user.UserClientMock{}
	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)
	client.On("FindOne", t.FindOneUserReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.FindOne(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *UserServiceTest) TestFindOneInternalError() {
	expected := t.InternalErr

	client := user.UserClientMock{}
	clientErr := status.Error(codes.Internal, constant.InternalErrorMessage)
	client.On("FindOne", t.FindOneUserReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.FindOne(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *UserServiceTest) TestUpdateSuccess() {
	protoResp := &proto.UpdateUserResponse{
		User: &proto.User{
			Id:        t.User.Id,
			Email:     t.User.Email,
			Firstname: t.User.Firstname,
			Lastname:  t.User.Lastname,
			Role:      t.User.Role,
		},
	}

	expected := &dto.UpdateUserResponse{
		Id:        t.User.Id,
		Email:     t.User.Email,
		Firstname: t.User.Firstname,
		Lastname:  t.User.Lastname,
	}

	client := user.UserClientMock{}
	client.On("Update", t.UpdateUserReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.Update(t.User.Id, t.UpdateUserDto)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *UserServiceTest) TestUpdateDuplicateEmail() {
	expected := t.ConflictErr

	client := user.UserClientMock{}
	clientErr := status.Error(codes.AlreadyExists, constant.DuplicateEmailMessage)
	client.On("Update", t.UpdateUserReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Update(t.User.Id, t.UpdateUserDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *UserServiceTest) TestUpdateUnavailableServiceError() {
	expected := t.UnavailableServiceErr

	client := user.UserClientMock{}
	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)
	client.On("Update", t.UpdateUserReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Update(t.User.Id, t.UpdateUserDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *UserServiceTest) TestUpdateInternalError() {
	expected := t.InternalErr

	client := user.UserClientMock{}
	clientErr := status.Error(codes.Internal, constant.InternalErrorMessage)
	client.On("Update", t.UpdateUserReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Update(t.User.Id, t.UpdateUserDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *UserServiceTest) TestDeleteSuccess() {
	protoResp := &proto.DeleteUserResponse{
		Success: true,
	}

	expected := &dto.DeleteUserResponse{
		Success: true,
	}

	client := user.UserClientMock{}
	client.On("Delete", t.DeleteUserReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.Delete(t.User.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *UserServiceTest) TestDeleteUnavailableServiceError() {
	expected := t.UnavailableServiceErr

	client := user.UserClientMock{}
	clientErr := status.Error(codes.Unavailable, constant.UnavailableServiceMessage)
	client.On("Delete", t.DeleteUserReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Delete(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *UserServiceTest) TestDeleteInternalError() {
	expected := t.InternalErr

	client := user.UserClientMock{}
	clientErr := status.Error(codes.Internal, constant.InternalErrorMessage)
	client.On("Delete", t.DeleteUserReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Delete(t.User.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}
