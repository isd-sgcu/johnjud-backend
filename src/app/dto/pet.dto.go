package dto

import (
	"github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
)

type PetDto struct {
	Type       string     `json:"type" validate:"required"`
	Species    string     `json:"species" validate:"required"`
	Name       string     `json:"name" validate:"required"`
	Birthdate  string     `json:"birthdate" validate:"required"`
	Gender     pet.Gender `json:"gender" validate:"required" example:"male"`
	Sterile    bool       `json:"sterile" validate:"required"`
	Vaccine    bool       `json:"vaccine" validate:"required"`
	Status     pet.Status `json:"status" validate:"required" example:"findhome"`
	Habit      string     `json:"habit" validate:"required"`
	Caption    string     `json:"caption"`
	Visible    bool       `json:"visible" validate:"required"`
	IsClubPet  bool       `json:"is_club_pet" validate:"required"`
	Background string     `json:"background"`
	Address    string     `json:"address"`
	Contact    string     `json:"contact"`
}
