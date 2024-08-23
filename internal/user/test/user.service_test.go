package user

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/isd-sgcu/johnjud-backend/internal/dto"
	"github.com/isd-sgcu/johnjud-backend/internal/model"
	"github.com/isd-sgcu/johnjud-backend/internal/user"
	mock "github.com/isd-sgcu/johnjud-backend/mocks/repository/user"
	"github.com/isd-sgcu/johnjud-backend/mocks/utils"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserServiceTest struct {
	suite.Suite
	User              *model.User
	UpdateUser        *model.User
	UserDto           *dto.User
	UserDtoNoPassword *dto.User
	HashedPassword    string
	UpdateUserReqMock *dto.UpdateUserRequest
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTest))
}

func (t *UserServiceTest) SetupTest() {
	t.User = &model.User{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Email:     faker.Email(),
		Password:  faker.Password(),
		Firstname: faker.Username(),
		Lastname:  faker.Username(),
		Role:      "user",
	}

	t.UserDto = &dto.User{
		Id:        t.User.ID.String(),
		Email:     t.User.Email,
		Firstname: t.User.Firstname,
		Lastname:  t.User.Lastname,
	}

	t.UserDtoNoPassword = &dto.User{
		Id:        t.User.ID.String(),
		Email:     t.User.Email,
		Firstname: t.User.Firstname,
		Lastname:  t.User.Lastname,
	}

	t.UpdateUserReqMock = &dto.UpdateUserRequest{
		Email:     t.User.Email,
		Password:  t.User.Password,
		Firstname: t.User.Firstname,
		Lastname:  t.User.Lastname,
	}

	t.HashedPassword = faker.Password()

	t.UpdateUser = &model.User{
		Email:     t.User.Email,
		Password:  t.HashedPassword,
		Firstname: t.User.Firstname,
		Lastname:  t.User.Lastname,
	}
}

func (t *UserServiceTest) TestFindOneSuccess() {
	want := t.UserDtoNoPassword

	repo := &mock.UserRepositoryMock{}
	repo.On("FindById", t.User.ID.String(), &model.User{}).Return(t.User, nil)

	brcyptUtil := &utils.BcryptUtilMock{}
	srv := user.NewService(repo, brcyptUtil)
	actual, err := srv.FindOne(t.User.ID.String())

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestFindOneNotFoundErr() {
	repo := &mock.UserRepositoryMock{}
	repo.On("FindById", t.User.ID.String(), &model.User{}).Return(nil, gorm.ErrRecordNotFound)

	brcyptUtil := &utils.BcryptUtilMock{}
	srv := user.NewService(repo, brcyptUtil)
	actual, err := srv.FindOne(t.User.ID.String())

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), http.StatusNotFound, err.StatusCode)
}

func (t *UserServiceTest) TestFindOneInternalErr() {
	repo := &mock.UserRepositoryMock{}
	repo.On("FindById", t.User.ID.String(), &model.User{}).Return(nil, errors.New("Not found user"))

	brcyptUtil := &utils.BcryptUtilMock{}
	srv := user.NewService(repo, brcyptUtil)
	actual, err := srv.FindOne(t.User.ID.String())

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), http.StatusInternalServerError, err.StatusCode)
}

func (t *UserServiceTest) TestUpdateSuccess() {
	want := t.UserDtoNoPassword

	repo := &mock.UserRepositoryMock{}
	repo.On("Update", t.User.ID.String(), t.UpdateUser).Return(t.User, nil)

	brcyptUtil := &utils.BcryptUtilMock{}
	brcyptUtil.On("GenerateHashedPassword", t.User.Password).Return(t.HashedPassword, nil)

	srv := user.NewService(repo, brcyptUtil)
	actual, err := srv.Update(t.User.ID.String(), t.UpdateUserReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestUpdateInternalErr() {
	repo := &mock.UserRepositoryMock{}
	repo.On("Update", t.User.ID.String(), t.UpdateUser).Return(nil, errors.New("Not found user"))

	brcyptUtil := &utils.BcryptUtilMock{}
	brcyptUtil.On("GenerateHashedPassword", t.User.Password).Return(t.HashedPassword, nil)

	srv := user.NewService(repo, brcyptUtil)
	actual, err := srv.Update(t.User.ID.String(), t.UpdateUserReqMock)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), http.StatusInternalServerError, err.StatusCode)
}

func (t *UserServiceTest) TestDeleteSuccess() {
	want := &dto.DeleteUserResponse{Success: true}

	repo := &mock.UserRepositoryMock{}
	repo.On("Delete", t.User.ID.String()).Return(nil)

	brcyptUtil := &utils.BcryptUtilMock{}
	srv := user.NewService(repo, brcyptUtil)
	actual, err := srv.Delete(t.UserDto.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestDeleteInternalErr() {
	repo := &mock.UserRepositoryMock{}
	repo.On("Delete", t.User.ID.String()).Return(errors.New("Not found user"))

	brcyptUtil := &utils.BcryptUtilMock{}
	srv := user.NewService(repo, brcyptUtil)
	actual, err := srv.Delete(t.UserDto.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), http.StatusInternalServerError, err.StatusCode)
}
