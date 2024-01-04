package auth

import (
	"errors"
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
	signupRequestDto    *dto.SignupRequest
	signInDto           *dto.SignInRequest
	token               string
	refreshTokenRequest *dto.RefreshTokenRequest
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
	signInDto := &dto.SignInRequest{
		Email:    faker.Email(),
		Password: faker.Password(),
	}
	token := faker.Word()
	refreshTokenRequest := &dto.RefreshTokenRequest{
		RefreshToken: faker.UUIDDigit(),
	}

	t.signupRequestDto = signupRequestDto
	t.signInDto = signInDto
	t.token = token
	t.refreshTokenRequest = refreshTokenRequest
}

func (t *AuthServiceTest) TestSignupSuccess() {
	protoReq := &authProto.SignUpRequest{
		FirstName: t.signupRequestDto.Firstname,
		LastName:  t.signupRequestDto.Lastname,
		Email:     t.signupRequestDto.Email,
		Password:  t.signupRequestDto.Password,
	}
	protoResp := &authProto.SignUpResponse{
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

	client.On("SignUp", protoReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.Signup(t.signupRequestDto)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected.Id, actual.Id)
	assert.Equal(t.T(), expected.Firstname, actual.Firstname)
	assert.Equal(t.T(), expected.Lastname, actual.Lastname)
	assert.Equal(t.T(), expected.Email, actual.Email)
}

func (t *AuthServiceTest) TestSignupConflict() {
	protoReq := &authProto.SignUpRequest{
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
	client.On("SignUp", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Signup(t.signupRequestDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestSignupInternalError() {
	protoReq := &authProto.SignUpRequest{
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
	client.On("SignUp", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Signup(t.signupRequestDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestSignupUnavailableService() {
	protoReq := &authProto.SignUpRequest{
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
	client.On("SignUp", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Signup(t.signupRequestDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestSignupUnknownError() {
	protoReq := &authProto.SignUpRequest{
		FirstName: t.signupRequestDto.Firstname,
		LastName:  t.signupRequestDto.Lastname,
		Email:     t.signupRequestDto.Email,
		Password:  t.signupRequestDto.Password,
	}
	clientErr := errors.New("Unknown error")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("SignUp", protoReq).Return(nil, clientErr)

	svc := NewService(&client)
	actual, err := svc.Signup(t.signupRequestDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestSignInSuccess() {
	protoReq := &authProto.SignInRequest{
		Email:    t.signInDto.Email,
		Password: t.signInDto.Password,
	}
	protoResp := &authProto.SignInResponse{
		Credential: &authProto.Credential{
			AccessToken:  faker.Word(),
			RefreshToken: faker.UUIDDigit(),
			ExpiresIn:    3600,
		},
	}

	expected := &dto.Credential{
		AccessToken:  protoResp.Credential.AccessToken,
		RefreshToken: protoResp.Credential.RefreshToken,
		ExpiresIn:    int(protoResp.Credential.ExpiresIn),
	}

	client := auth.AuthClientMock{}

	client.On("SignIn", protoReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.SignIn(t.signInDto)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *AuthServiceTest) TestSignInForbidden() {
	protoReq := &authProto.SignInRequest{
		Email:    t.signInDto.Email,
		Password: t.signInDto.Password,
	}
	protoErr := status.Error(codes.PermissionDenied, "Incorrect email or password")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusForbidden,
		Message:    constant.IncorrectEmailPasswordMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}

	client.On("SignIn", protoReq).Return(nil, protoErr)

	svc := NewService(&client)
	actual, err := svc.SignIn(t.signInDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestSignInInternalError() {
	protoReq := &authProto.SignInRequest{
		Email:    t.signInDto.Email,
		Password: t.signInDto.Password,
	}
	protoErr := status.Error(codes.Internal, "Internal error")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}

	client.On("SignIn", protoReq).Return(nil, protoErr)

	svc := NewService(&client)
	actual, err := svc.SignIn(t.signInDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestSignInUnavailableService() {
	protoReq := &authProto.SignInRequest{
		Email:    t.signInDto.Email,
		Password: t.signInDto.Password,
	}
	protoErr := status.Error(codes.Unavailable, "Connection lost")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    constant.UnavailableServiceMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("SignIn", protoReq).Return(nil, protoErr)

	svc := NewService(&client)
	actual, err := svc.SignIn(t.signInDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestSignInUnknownError() {
	protoReq := &authProto.SignInRequest{
		Email:    t.signInDto.Email,
		Password: t.signInDto.Password,
	}
	protoErr := errors.New("Unknown error")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}

	client.On("SignIn", protoReq).Return(nil, protoErr)

	svc := NewService(&client)
	actual, err := svc.SignIn(t.signInDto)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestSignOutSuccess() {
	protoReq := &authProto.SignOutRequest{
		Token: t.token,
	}
	protoResp := &authProto.SignOutResponse{
		IsSuccess: true,
	}

	expected := &dto.SignOutResponse{IsSuccess: true}

	client := auth.AuthClientMock{}
	client.On("SignOut", protoReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.SignOut(t.token)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *AuthServiceTest) TestSignOutInternalError() {
	protoReq := &authProto.SignOutRequest{
		Token: t.token,
	}
	protoErr := status.Error(codes.Internal, "Internal error")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("SignOut", protoReq).Return(nil, protoErr)

	svc := NewService(&client)
	actual, err := svc.SignOut(t.token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestSignOutUnavailableService() {
	protoReq := &authProto.SignOutRequest{
		Token: t.token,
	}
	protoErr := status.Error(codes.Unavailable, "Connection lost")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    constant.UnavailableServiceMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("SignOut", protoReq).Return(nil, protoErr)

	svc := NewService(&client)
	actual, err := svc.SignOut(t.token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestSignOutUnknownError() {
	protoReq := &authProto.SignOutRequest{
		Token: t.token,
	}
	protoErr := errors.New("Unknown error")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("SignOut", protoReq).Return(nil, protoErr)

	svc := NewService(&client)
	actual, err := svc.SignOut(t.token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestValidateSuccess() {
	protoReq := &authProto.ValidateRequest{
		Token: t.token,
	}
	protoResp := &authProto.ValidateResponse{
		UserId: faker.UUIDDigit(),
		Role:   "user",
	}

	expected := &dto.TokenPayloadAuth{
		UserId: protoResp.UserId,
		Role:   protoResp.Role,
	}

	client := auth.AuthClientMock{}
	client.On("Validate", protoReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.Validate(t.token)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *AuthServiceTest) TestValidateUnauthorized() {
	protoReq := &authProto.ValidateRequest{
		Token: t.token,
	}
	protoErr := status.Error(codes.Unauthenticated, "invalid token")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusUnauthorized,
		Message:    constant.UnauthorizedMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("Validate", protoReq).Return(nil, protoErr)

	svc := NewService(&client)
	actual, err := svc.Validate(t.token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestValidateUnavailableService() {
	protoReq := &authProto.ValidateRequest{
		Token: t.token,
	}
	protoErr := status.Error(codes.Unavailable, "connection lost")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    constant.UnavailableServiceMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("Validate", protoReq).Return(nil, protoErr)

	svc := NewService(&client)
	actual, err := svc.Validate(t.token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestValidateUnknownError() {
	protoReq := &authProto.ValidateRequest{
		Token: t.token,
	}
	protoErr := errors.New("Unknown error")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("Validate", protoReq).Return(nil, protoErr)

	svc := NewService(&client)
	actual, err := svc.Validate(t.token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestRefreshTokenSuccess() {
	protoReq := &authProto.RefreshTokenRequest{
		RefreshToken: t.refreshTokenRequest.RefreshToken,
	}
	protoResp := &authProto.RefreshTokenResponse{
		Credential: &authProto.Credential{
			AccessToken:  faker.Word(),
			RefreshToken: faker.UUIDDigit(),
			ExpiresIn:    3600,
		},
	}

	expected := &dto.Credential{
		AccessToken:  protoResp.Credential.AccessToken,
		RefreshToken: protoResp.Credential.RefreshToken,
		ExpiresIn:    int(protoResp.Credential.ExpiresIn),
	}

	client := auth.AuthClientMock{}
	client.On("RefreshToken", protoReq).Return(protoResp, nil)

	svc := NewService(&client)
	actual, err := svc.RefreshToken(t.refreshTokenRequest)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *AuthServiceTest) TestRefreshTokenInvalidToken() {
	protoReq := &authProto.RefreshTokenRequest{
		RefreshToken: t.refreshTokenRequest.RefreshToken,
	}
	protoErr := status.Error(codes.InvalidArgument, "Invalid token")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    constant.InvalidTokenMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("RefreshToken", protoReq).Return(nil, protoErr)

	svc := NewService(&client)
	actual, err := svc.RefreshToken(t.refreshTokenRequest)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestRefreshTokenInternalError() {
	protoReq := &authProto.RefreshTokenRequest{
		RefreshToken: t.refreshTokenRequest.RefreshToken,
	}
	protoErr := status.Error(codes.Internal, "Internal error")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("RefreshToken", protoReq).Return(nil, protoErr)

	svc := NewService(&client)
	actual, err := svc.RefreshToken(t.refreshTokenRequest)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestRefreshTokenUnavailableService() {
	protoReq := &authProto.RefreshTokenRequest{
		RefreshToken: t.refreshTokenRequest.RefreshToken,
	}
	protoErr := status.Error(codes.Unavailable, "Connection lost")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    constant.UnavailableServiceMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("RefreshToken", protoReq).Return(nil, protoErr)

	svc := NewService(&client)
	actual, err := svc.RefreshToken(t.refreshTokenRequest)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *AuthServiceTest) TestRefreshTokenUnknownError() {
	protoReq := &authProto.RefreshTokenRequest{
		RefreshToken: t.refreshTokenRequest.RefreshToken,
	}
	protoErr := errors.New("Unknown error")

	expected := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    constant.InternalErrorMessage,
		Data:       nil,
	}

	client := auth.AuthClientMock{}
	client.On("RefreshToken", protoReq).Return(nil, protoErr)

	svc := NewService(&client)
	actual, err := svc.RefreshToken(t.refreshTokenRequest)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}
