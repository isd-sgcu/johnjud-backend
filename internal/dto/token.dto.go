package dto

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/isd-sgcu/johnjud-gateway/constant"
)

type UserCredential struct {
	UserID        string        `json:"user_id"`
	Role          constant.Role `json:"role"`
	AuthSessionID string        `json:"auth_session_id"`
	RefreshToken  string        `json:"refresh_token"`
}

type AuthPayload struct {
	jwt.RegisteredClaims
	UserID        string `json:"user_id"`
	AuthSessionID string `json:"auth_session_id"`
}

type AccessTokenCache struct {
	Token        string        `json:"token"`
	Role         constant.Role `json:"role"`
	RefreshToken string        `json:"refresh_token"`
}

type RefreshTokenCache struct {
	AuthSessionID string        `json:"auth_session_id"`
	UserID        string        `json:"user_id"`
	Role          constant.Role `json:"role"`
}

type ResetPasswordTokenCache struct {
	UserID string `json:"user_id"`
}
