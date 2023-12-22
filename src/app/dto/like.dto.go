package dto

import (
	"github.com/google/uuid"
)

type LikeDto struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	PetID  uuid.UUID `json:"pet_id" validate:"required"`
}
