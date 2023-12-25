package auth

import (
	"github.com/go-faker/faker/v4"
	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/mocks/client/auth"
	authProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/auth/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"testing"
)

type AuthServiceTest struct {
	suite.Suite
	signupRequestDto *dto.SignupRequest
	signInDto        *dto.SignIn
	token            string
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTest))
}

func (t *AuthServiceTest) SetupTest() {
	signupRequestDto := &dto.SignupRequest{
		Email:     faker.Email(),
		Password:  faker.Password(),
		Firstname: faker.FirstName(),
		Lastname:  faker.LastName(),
	}
	signInDto := &dto.SignIn{
		Email:    faker.Email(),
		Password: faker.Password(),
	}
	token := faker.Word()

	t.signupRequestDto = signupRequestDto
	t.signInDto = signInDto
	t.token = token
}

func (t *AuthServiceTest) TestSignupSuccess() {
	protoReq := &authProto.SignupRequest{
		FirstName: t.signupRequestDto.Firstname,
		LastName:  t.signupRequestDto.Lastname,
		Email:     t.signupRequestDto.Email,
		Password:  t.signupRequestDto.Password,
	}
	protoResp := &authProto.SignupResponse{
		Id:        faker.UUIDDigit(),
		FirstName: t.signupRequestDto.Firstname,
		LastName:  t.signupRequestDto.Lastname,
		Email:     t.signupRequestDto.Email,
	}

	expected := &dto.SignupResponse{
		Id:        protoResp.Id,
		Email:     t.signupRequestDto.Email,
		Firstname: t.signupRequestDto.Firstname,
		Lastname:  t.signupRequestDto.Lastname,
	}

	client := auth.AuthClientMock{}

	client.On("Signup", protoReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.Signup(t.signupRequestDto)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected.Id, actual.Id)
	assert.Equal(t.T(), expected.Firstname, actual.Firstname)
	assert.Equal(t.T(), expected.Lastname, actual.Lastname)
	assert.Equal(t.T(), expected.Email, actual.Email)
}

func (t *AuthServiceTest) TestSignupConflict() {
	protoReq := &authProto.SignupRequest{
		FirstName: t.signupRequestDto.Firstname,
		LastName:  t.signupRequestDto.Lastname,
		Email:     t.signupRequestDto.Email,
		Password:  t.signupRequestDto.Password,
	}
	clientErr := status.Error(codes.AlreadyExists, "Duplicate email")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusConflict,
		Message:    constant.DuplicateEmailMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("Signup", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Signup(t.signupRequestDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestSignupInternalError() {
	protoReq := &authProto.SignupRequest{
		FirstName: t.signupRequestDto.Firstname,
		LastName:  t.signupRequestDto.Lastname,
		Email:     t.signupRequestDto.Email,
		Password:  t.signupRequestDto.Password,
	}
	clientErr := status.Error(codes.Internal, "Internal error")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("Signup", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Signup(t.signupRequestDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestSignupUnavailableService() {
	protoReq := &authProto.SignupRequest{
		FirstName: t.signupRequestDto.Firstname,
		LastName:  t.signupRequestDto.Lastname,
		Email:     t.signupRequestDto.Email,
		Password:  t.signupRequestDto.Password,
	}
	clientErr := status.Error(codes.Unavailable, "Connection lost")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    constant.UnavailableServiceMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("Signup", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Signup(t.signupRequestDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestSignInSuccess() {}

func (t *AuthServiceTest) TestSignInForbidden() {}

func (t *AuthServiceTest) TestSignInInternalError() {}

func (t *AuthServiceTest) TestSignInUnavailableService() {}

func (t *AuthServiceTest) TestValidateSuccess() {}

func (t *AuthServiceTest) TestValidateUnauthorized() {}

func (t *AuthServiceTest) TestValidateInternalError() {}

func (t *AuthServiceTest) TestValidateUnavailableService() {}

func (t *AuthServiceTest) TestRefreshTokenUnauthorized() {}

func (t *AuthServiceTest) TestRefreshTokenInternalError() {}

func (t *AuthServiceTest) TestRefreshTokenUnavailableService() {}
