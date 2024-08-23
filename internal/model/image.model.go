package model

import (
	"github.com/google/uuid"
)

type Image struct {
	Base
	PetID     *uuid.UUID `json:"pet_id" gorm:"index:idx_name"`
	Pet       *Pet       `json:"pet" gorm:"foreignKey:PetID;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
	ImageUrl  string     `json:"image_url" gorm:"mediumtext"`
	ObjectKey string     `json:"object_key" gorm:"mediumtext"`
}
