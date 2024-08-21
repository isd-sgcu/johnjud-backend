package model

import "github.com/google/uuid"

type AuthSession struct {
	Base
	UserID uuid.UUID `json:"user_id"`
}
