package token

import (
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	_jwt "github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-gateway/config"
	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/auth/token"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	mock_cache "github.com/isd-sgcu/johnjud-gateway/mocks/repository/cache"
	"github.com/isd-sgcu/johnjud-gateway/mocks/service/jwt"
	"github.com/isd-sgcu/johnjud-gateway/mocks/utils"
	authProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/auth/v1"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenServiceTest struct {
	suite.Suite
	userId        string
	role          constant.Role
	authSessionId string
	accessToken   string
	refreshToken  *uuid.UUID
	jwtConfig     *config.Jwt
	validateToken string
}

func TestTokenService(t *testing.T) {
	suite.Run(t, new(TokenServiceTest))
}

func (t *TokenServiceTest) SetupTest() {
	userId := faker.UUIDDigit()
	role := constant.USER
	authSessionId := faker.UUIDDigit()
	accessToken := "testAccessToken"
	refreshToken := uuid.New()
	jwtConfig := &config.Jwt{
		Secret:          "testSecret",
		ExpiresIn:       3600,
		RefreshTokenTTL: 604800,
		Issuer:          "testIssuer",
		ResetTokenTTL:   900,
	}
	validateToken := ""

	t.userId = userId
	t.role = role
	t.authSessionId = authSessionId
	t.accessToken = accessToken
	t.refreshToken = &refreshToken
	t.jwtConfig = jwtConfig
	t.validateToken = validateToken
}

func (t *TokenServiceTest) TestCreateCredentialSuccess() {
	accessTokenCache := &dto.AccessTokenCache{
		Token:        t.accessToken,
		Role:         t.role,
		RefreshToken: t.refreshToken.String(),
	}
	refreshTokenCache := &dto.RefreshTokenCache{
		AuthSessionID: t.authSessionId,
		UserID:        t.userId,
		Role:          t.role,
	}

	expected := authProto.Credential{
		AccessToken:  t.accessToken,
		RefreshToken: t.refreshToken.String(),
		ExpiresIn:    int32(t.jwtConfig.ExpiresIn),
	}

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	jwtService.On("SignAuth", t.userId, t.role, t.authSessionId).Return(t.accessToken, nil)
	jwtService.On("GetConfig").Return(t.jwtConfig)
	uuidUtil.On("GetNewUUID").Return(t.refreshToken)
	accessTokenRepo.EXPECT().SetValue(t.authSessionId, accessTokenCache, t.jwtConfig.ExpiresIn).Return(nil)
	refreshTokenRepo.EXPECT().SetValue(t.refreshToken.String(), refreshTokenCache, t.jwtConfig.RefreshTokenTTL).Return(nil)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.CreateCredential(t.userId, t.role, t.authSessionId)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected.AccessToken, actual.AccessToken)
	assert.Equal(t.T(), expected.RefreshToken, actual.RefreshToken)
	assert.Equal(t.T(), expected.ExpiresIn, actual.ExpiresIn)
}

func (t *TokenServiceTest) TestCreateCredentialSignAuthFailed() {
	signAuthError := errors.New("Error while signing token")
	expected := errors.New("Error while signing token")

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	jwtService.On("SignAuth", t.userId, t.role, t.authSessionId).Return("", signAuthError)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.CreateCredential(t.userId, t.role, t.authSessionId)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected.Error(), err.Error())
}

func (t *TokenServiceTest) TestCreateCredentialSetAccessTokenFailed() {
	accessTokenCache := &dto.AccessTokenCache{
		Token:        t.accessToken,
		Role:         t.role,
		RefreshToken: t.refreshToken.String(),
	}
	setCacheErr := errors.New("Internal server error")
	expected := setCacheErr

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	jwtService.On("SignAuth", t.userId, t.role, t.authSessionId).Return(t.accessToken, nil)
	jwtService.On("GetConfig").Return(t.jwtConfig)
	uuidUtil.On("GetNewUUID").Return(t.refreshToken)
	accessTokenRepo.EXPECT().SetValue(t.authSessionId, accessTokenCache, t.jwtConfig.ExpiresIn).Return(setCacheErr)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.CreateCredential(t.userId, t.role, t.authSessionId)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected.Error(), err.Error())
}

func (t *TokenServiceTest) TestCreateCredentialSetRefreshTokenFailed() {
	accessTokenCache := &dto.AccessTokenCache{
		Token:        t.accessToken,
		Role:         t.role,
		RefreshToken: t.refreshToken.String(),
	}
	refreshTokenCache := &dto.RefreshTokenCache{
		AuthSessionID: t.authSessionId,
		UserID:        t.userId,
		Role:          t.role,
	}
	setCacheErr := errors.New("Internal server error")
	expected := setCacheErr

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	jwtService.On("SignAuth", t.userId, t.role, t.authSessionId).Return(t.accessToken, nil)
	jwtService.On("GetConfig").Return(t.jwtConfig)
	uuidUtil.On("GetNewUUID").Return(t.refreshToken)
	accessTokenRepo.EXPECT().SetValue(t.authSessionId, accessTokenCache, t.jwtConfig.ExpiresIn).Return(nil)
	refreshTokenRepo.EXPECT().SetValue(t.refreshToken.String(), refreshTokenCache, t.jwtConfig.RefreshTokenTTL).Return(setCacheErr)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.CreateCredential(t.userId, t.role, t.authSessionId)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected.Error(), err.Error())
}

func (t *TokenServiceTest) TestValidateSuccess() {
	expected := &dto.UserCredential{
		UserID:        t.userId,
		Role:          "",
		AuthSessionID: t.authSessionId,
		RefreshToken:  "",
	}
	payloads := _jwt.MapClaims{
		"iss":             t.jwtConfig.Issuer,
		"exp":             float64(_jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(t.jwtConfig.ExpiresIn))).Unix()),
		"iat":             float64(_jwt.NewNumericDate(time.Now()).Unix()),
		"user_id":         t.userId,
		"auth_session_id": t.authSessionId,
	}
	jwtToken := &_jwt.Token{
		Method: _jwt.SigningMethodHS256,
		Claims: payloads,
	}
	accessTokenCache := &dto.AccessTokenCache{}

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	jwtService.On("VerifyAuth", t.validateToken).Return(jwtToken, nil)
	jwtService.On("GetConfig").Return(t.jwtConfig)
	accessTokenRepo.EXPECT().GetValue(payloads["auth_session_id"].(string), accessTokenCache).Return(nil)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.Validate(t.validateToken)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), *expected, *actual)
}

func (t *TokenServiceTest) TestValidateInvalidIssuer() {
	expected := errors.New("invalid token")

	payloads := _jwt.MapClaims{
		"iss":             "invalid issuer",
		"exp":             float64(_jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(t.jwtConfig.ExpiresIn))).Unix()),
		"iat":             float64(_jwt.NewNumericDate(time.Now()).Unix()),
		"user_id":         t.userId,
		"auth_session_id": t.authSessionId,
	}

	jwtToken := &_jwt.Token{
		Method: _jwt.SigningMethodHS256,
		Claims: payloads,
	}

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	jwtService.On("VerifyAuth", t.validateToken).Return(jwtToken, nil)
	jwtService.On("GetConfig").Return(t.jwtConfig)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.Validate(t.validateToken)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected.Error(), err.Error())
}

func (t *TokenServiceTest) TestValidateExpireToken() {
	expected := errors.New("expired token")

	payloads := _jwt.MapClaims{
		"iss":             t.jwtConfig.Issuer,
		"exp":             float64(_jwt.NewNumericDate(time.Now().Add(time.Second * (-time.Duration(t.jwtConfig.ExpiresIn)))).Unix()),
		"iat":             float64(_jwt.NewNumericDate(time.Now()).Unix()),
		"user_id":         t.userId,
		"auth_session_id": t.authSessionId,
	}
	jwtToken := &_jwt.Token{
		Method: _jwt.SigningMethodHS256,
		Claims: payloads,
	}

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	jwtService.On("VerifyAuth", t.validateToken).Return(jwtToken, nil)
	jwtService.On("GetConfig").Return(t.jwtConfig)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.Validate(t.validateToken)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected.Error(), err.Error())
}

func (t *TokenServiceTest) TestValidateVerifyFailed() {
	expected := errors.New("invalid token")

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	jwtService.On("VerifyAuth", t.validateToken).Return(nil, expected)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.Validate(t.validateToken)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected.Error(), err.Error())
}

func (t *TokenServiceTest) TestValidateGetCacheKeyNotFound() {
	expected := errors.New("invalid token")

	payloads := _jwt.MapClaims{
		"iss":             t.jwtConfig.Issuer,
		"exp":             float64(_jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(t.jwtConfig.ExpiresIn))).Unix()),
		"iat":             float64(_jwt.NewNumericDate(time.Now()).Unix()),
		"user_id":         t.userId,
		"auth_session_id": t.authSessionId,
	}
	jwtToken := &_jwt.Token{
		Method: _jwt.SigningMethodHS256,
		Claims: payloads,
	}
	accessTokenCache := &dto.AccessTokenCache{}

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	jwtService.On("VerifyAuth", t.validateToken).Return(jwtToken, nil)
	jwtService.On("GetConfig").Return(t.jwtConfig)
	accessTokenRepo.EXPECT().GetValue(payloads["auth_session_id"].(string), accessTokenCache).Return(redis.Nil)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.Validate(t.validateToken)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected.Error(), err.Error())
}

func (t *TokenServiceTest) TestValidateGetCacheInternalFailed() {
	payloads := _jwt.MapClaims{
		"iss":             t.jwtConfig.Issuer,
		"exp":             float64(_jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(t.jwtConfig.ExpiresIn))).Unix()),
		"iat":             float64(_jwt.NewNumericDate(time.Now()).Unix()),
		"user_id":         t.userId,
		"auth_session_id": t.authSessionId,
	}
	jwtToken := &_jwt.Token{
		Method: _jwt.SigningMethodHS256,
		Claims: payloads,
	}
	accessTokenCache := &dto.AccessTokenCache{}
	getCacheErr := errors.New("internal server error")

	expected := getCacheErr

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	jwtService.On("VerifyAuth", t.validateToken).Return(jwtToken, nil)
	jwtService.On("GetConfig").Return(t.jwtConfig)
	accessTokenRepo.EXPECT().GetValue(payloads["auth_session_id"].(string), accessTokenCache).Return(getCacheErr)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.Validate(t.validateToken)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected.Error(), err.Error())
}

func (t *TokenServiceTest) TestValidateInvalidToken() {
	invalidToken := faker.Word()
	expected := errors.New("invalid token")

	payloads := _jwt.MapClaims{
		"iss":             t.jwtConfig.Issuer,
		"exp":             float64(_jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(t.jwtConfig.ExpiresIn))).Unix()),
		"iat":             float64(_jwt.NewNumericDate(time.Now()).Unix()),
		"user_id":         t.userId,
		"auth_session_id": t.authSessionId,
	}
	jwtToken := &_jwt.Token{
		Method: _jwt.SigningMethodHS256,
		Claims: payloads,
	}
	accessTokenCache := &dto.AccessTokenCache{}

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	jwtService.On("VerifyAuth", invalidToken).Return(jwtToken, nil)
	jwtService.On("GetConfig").Return(t.jwtConfig)
	accessTokenRepo.EXPECT().GetValue(payloads["auth_session_id"].(string), accessTokenCache).Return(nil)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.Validate(invalidToken)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected.Error(), err.Error())
}

func (t *TokenServiceTest) TestCreateRefreshTokenSuccess() {
	expected := t.refreshToken.String()

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	uuidUtil.On("GetNewUUID").Return(t.refreshToken)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual := tokenSvc.CreateRefreshToken()

	assert.Equal(t.T(), expected, actual)
}

func (t *TokenServiceTest) TestRemoveAccessTokenCacheSuccess() {
	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	accessTokenRepo.EXPECT().DeleteValue(t.authSessionId).Return(nil)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	err := tokenSvc.RemoveAccessTokenCache(t.authSessionId)

	assert.Nil(t.T(), err)
}

func (t *TokenServiceTest) TestRemoveAccessTokenCacheDeleteInternalFailed() {
	deleteAccessTokenCacheErr := errors.New("internal server error")

	expected := deleteAccessTokenCacheErr

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	accessTokenRepo.EXPECT().DeleteValue(t.authSessionId).Return(deleteAccessTokenCacheErr)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	err := tokenSvc.RemoveAccessTokenCache(t.authSessionId)

	assert.Equal(t.T(), expected, err)
}

func (t *TokenServiceTest) TestFindRefreshTokenCacheSuccess() {
	expected := &dto.RefreshTokenCache{}

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	refreshTokenRepo.EXPECT().GetValue(t.refreshToken.String(), &dto.RefreshTokenCache{}).Return(nil)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.FindRefreshTokenCache(t.refreshToken.String())

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *TokenServiceTest) TestFindRefreshTokenCacheInvalid() {
	getCacheErr := redis.Nil

	expected := status.Error(codes.InvalidArgument, getCacheErr.Error())

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	refreshTokenRepo.EXPECT().GetValue(t.refreshToken.String(), &dto.RefreshTokenCache{}).Return(getCacheErr)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.FindRefreshTokenCache(t.refreshToken.String())

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *TokenServiceTest) TestFindRefreshTokenCacheInternalError() {
	getCacheErr := errors.New("internal server error")

	expected := status.Error(codes.Internal, getCacheErr.Error())

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	refreshTokenRepo.EXPECT().GetValue(t.refreshToken.String(), &dto.RefreshTokenCache{}).Return(getCacheErr)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.FindRefreshTokenCache(t.refreshToken.String())

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *TokenServiceTest) TestRemoveRefreshTokenCacheSuccess() {
	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	refreshTokenRepo.EXPECT().DeleteValue(t.refreshToken.String()).Return(nil)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	err := tokenSvc.RemoveRefreshTokenCache(t.refreshToken.String())

	assert.Nil(t.T(), err)
}

func (t *TokenServiceTest) TestRemoveRefreshTokenCacheDeleteInternalFailed() {
	deleteRefreshTokenCacheErr := errors.New("internal server error")

	expected := deleteRefreshTokenCacheErr

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	refreshTokenRepo.EXPECT().DeleteValue(t.refreshToken.String()).Return(deleteRefreshTokenCacheErr)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	err := tokenSvc.RemoveRefreshTokenCache(t.refreshToken.String())

	assert.Equal(t.T(), expected, err)
}

func (t *TokenServiceTest) TestCreateResetPasswordTokenSuccess() {
	tokenCache := &dto.ResetPasswordTokenCache{
		UserID: t.userId,
	}

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	uuidUtil.On("GetNewUUID").Return(t.refreshToken)
	jwtService.On("GetConfig").Return(t.jwtConfig)
	resetPasswordTokenRepo.EXPECT().SetValue(t.refreshToken.String(), tokenCache, t.jwtConfig.ResetTokenTTL).Return(nil)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.CreateResetPasswordToken(t.userId)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), t.refreshToken.String(), actual)
}

func (t *TokenServiceTest) TestCreateResetPasswordTokenFailed() {
	tokenCache := &dto.ResetPasswordTokenCache{
		UserID: t.userId,
	}
	cacheErr := errors.New("Internal error")

	expected := cacheErr

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	uuidUtil.On("GetNewUUID").Return(t.refreshToken)
	jwtService.On("GetConfig").Return(t.jwtConfig)
	resetPasswordTokenRepo.EXPECT().SetValue(t.refreshToken.String(), tokenCache, t.jwtConfig.ResetTokenTTL).Return(cacheErr)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.CreateResetPasswordToken(t.userId)

	assert.Equal(t.T(), "", actual)
	assert.Equal(t.T(), expected, err)
}

func (t *TokenServiceTest) TestFindResetPasswordTokenSuccess() {
	tokenCache := &dto.ResetPasswordTokenCache{}

	expected := tokenCache

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	resetPasswordTokenRepo.EXPECT().GetValue(t.refreshToken.String(), tokenCache).Return(nil)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.FindResetPasswordToken(t.refreshToken.String())

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *TokenServiceTest) TestFindResetPasswordTokenNotFound() {
	tokenCache := &dto.ResetPasswordTokenCache{}
	cacheErr := redis.Nil

	expected := status.Error(codes.InvalidArgument, cacheErr.Error())

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	resetPasswordTokenRepo.EXPECT().GetValue(t.refreshToken.String(), tokenCache).Return(cacheErr)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.FindResetPasswordToken(t.refreshToken.String())

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *TokenServiceTest) TestFindResetPasswordTokenInternalError() {
	tokenCache := &dto.ResetPasswordTokenCache{}
	cacheErr := errors.New("Internal error")

	expected := status.Error(codes.Internal, cacheErr.Error())

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	resetPasswordTokenRepo.EXPECT().GetValue(t.refreshToken.String(), tokenCache).Return(cacheErr)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	actual, err := tokenSvc.FindResetPasswordToken(t.refreshToken.String())

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected, err)
}

func (t *TokenServiceTest) TestRemoveResetPasswordTokenSuccess() {
	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	resetPasswordTokenRepo.EXPECT().DeleteValue(t.refreshToken.String()).Return(nil)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	err := tokenSvc.RemoveResetPasswordToken(t.refreshToken.String())

	assert.Nil(t.T(), err)
}

func (t *TokenServiceTest) TestRemoveResetPasswordTokenFailed() {
	cacheErr := errors.New("Internal error")

	expected := cacheErr

	controller := gomock.NewController(t.T())

	jwtService := jwt.JwtServiceMock{}
	accessTokenRepo := mock_cache.NewMockRepository(controller)
	refreshTokenRepo := mock_cache.NewMockRepository(controller)
	resetPasswordTokenRepo := mock_cache.NewMockRepository(controller)
	uuidUtil := utils.UuidUtilMock{}

	resetPasswordTokenRepo.EXPECT().DeleteValue(t.refreshToken.String()).Return(cacheErr)

	tokenSvc := token.NewService(&jwtService, accessTokenRepo, refreshTokenRepo, resetPasswordTokenRepo, &uuidUtil)
	err := tokenSvc.RemoveResetPasswordToken(t.refreshToken.String())

	assert.Equal(t.T(), expected, err)
}
