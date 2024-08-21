package token

import (
	"time"

	_jwt "github.com/golang-jwt/jwt/v4"
	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/auth/jwt"
	"github.com/isd-sgcu/johnjud-gateway/internal/cache"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	"github.com/isd-sgcu/johnjud-gateway/internal/utils"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	CreateCredential(userId string, role constant.Role, authSessionId string) (*dto.Credential, error)
	Validate(token string) (*dto.UserCredential, error)
	CreateRefreshToken() string
	RemoveAccessTokenCache(authSessionId string) error
	FindRefreshTokenCache(refreshToken string) (*dto.RefreshTokenCache, error)
	RemoveRefreshTokenCache(refreshToken string) error
	CreateResetPasswordToken(userId string) (string, error)
	FindResetPasswordToken(token string) (*dto.ResetPasswordTokenCache, error)
	RemoveResetPasswordToken(token string) error
}

type serviceImpl struct {
	jwtService              jwt.Service
	accessTokenCache        cache.Repository
	refreshTokenCache       cache.Repository
	resetPasswordTokenCache cache.Repository
	uuidUtil                utils.IUuidUtil
}

func NewService(jwtService jwt.Service, accessTokenCache cache.Repository, refreshTokenCache cache.Repository, resetPasswordTokenCache cache.Repository, uuidUtil utils.IUuidUtil) Service {
	return &serviceImpl{
		jwtService:              jwtService,
		accessTokenCache:        accessTokenCache,
		refreshTokenCache:       refreshTokenCache,
		resetPasswordTokenCache: resetPasswordTokenCache,
		uuidUtil:                uuidUtil,
	}
}

func (s *serviceImpl) CreateCredential(userId string, role constant.Role, authSessionId string) (*dto.Credential, error) {
	accessToken, err := s.jwtService.SignAuth(userId, role, authSessionId)
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "token").
			Str("module", "CreateCredential").
			Msg("Error signing jwt access token")
		return nil, err
	}

	refreshToken := s.CreateRefreshToken()
	jwtConf := s.jwtService.GetConfig()

	accessTokenCache := &dto.AccessTokenCache{
		Token:        accessToken,
		Role:         role,
		RefreshToken: refreshToken,
	}
	err = s.accessTokenCache.SetValue(authSessionId, accessTokenCache, jwtConf.ExpiresIn)
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "token").
			Str("module", "CreateCredential").
			Msg("Error setting value to access token cache")
		return nil, err
	}

	refreshTokenCache := &dto.RefreshTokenCache{
		AuthSessionID: authSessionId,
		UserID:        userId,
		Role:          role,
	}
	err = s.refreshTokenCache.SetValue(refreshToken, refreshTokenCache, jwtConf.RefreshTokenTTL)
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "token").
			Str("module", "CreateCredential").
			Msg("Error setting value to refresh token cache")
		return nil, err
	}

	credential := &dto.Credential{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    jwtConf.ExpiresIn,
	}

	return credential, nil
}

func (s *serviceImpl) Validate(token string) (*dto.UserCredential, error) {
	jwtToken, err := s.jwtService.VerifyAuth(token)
	if err != nil {
		return nil, err
	}

	payloads := jwtToken.Claims.(_jwt.MapClaims)
	if payloads["iss"] != s.jwtService.GetConfig().Issuer {
		return nil, errors.New("invalid token")
	}

	if time.Unix(int64(payloads["exp"].(float64)), 0).Before(time.Now()) {
		return nil, errors.New("expired token")
	}

	accessTokenCache := &dto.AccessTokenCache{}
	err = s.accessTokenCache.GetValue(payloads["auth_session_id"].(string), accessTokenCache)
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
		return nil, errors.New("invalid token")
	}

	if token != accessTokenCache.Token {
		return nil, errors.New("invalid token")
	}

	userCredential := &dto.UserCredential{
		UserID:        payloads["user_id"].(string),
		Role:          accessTokenCache.Role,
		AuthSessionID: payloads["auth_session_id"].(string),
		RefreshToken:  accessTokenCache.RefreshToken,
	}
	return userCredential, nil
}

func (s *serviceImpl) CreateRefreshToken() string {
	return s.uuidUtil.GetNewUUID().String()
}

func (s *serviceImpl) RemoveAccessTokenCache(authSessionId string) error {
	err := s.accessTokenCache.DeleteValue(authSessionId)
	if err != nil {
		if err != redis.Nil {
			return err
		}
	}

	return nil
}

func (s *serviceImpl) FindRefreshTokenCache(refreshToken string) (*dto.RefreshTokenCache, error) {
	refreshTokenCache := &dto.RefreshTokenCache{}
	err := s.refreshTokenCache.GetValue(refreshToken, refreshTokenCache)
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "token").
			Str("module", "FindRefreshTokenCache").
			Msg("Error getting value from redis")
		if err != redis.Nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return refreshTokenCache, nil
}

func (s *serviceImpl) RemoveRefreshTokenCache(refreshToken string) error {
	err := s.refreshTokenCache.DeleteValue(refreshToken)
	if err != nil {
		if err != redis.Nil {
			return err
		}
	}

	return nil
}

func (s *serviceImpl) CreateResetPasswordToken(userId string) (string, error) {
	resetPasswordToken := s.CreateRefreshToken()
	tokenCache := &dto.ResetPasswordTokenCache{
		UserID: userId,
	}
	err := s.resetPasswordTokenCache.SetValue(resetPasswordToken, tokenCache, s.jwtService.GetConfig().ResetTokenTTL)
	if err != nil {
		return "", err
	}
	return resetPasswordToken, nil
}

func (s *serviceImpl) FindResetPasswordToken(token string) (*dto.ResetPasswordTokenCache, error) {
	tokenCache := &dto.ResetPasswordTokenCache{}
	err := s.resetPasswordTokenCache.GetValue(token, tokenCache)
	if err != nil {
		if err != redis.Nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return tokenCache, nil
}

func (s *serviceImpl) RemoveResetPasswordToken(token string) error {
	err := s.resetPasswordTokenCache.DeleteValue(token)
	if err != nil {
		if err != redis.Nil {
			return err
		}
	}

	return nil
}
