package dto

import (
	"github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
)

type PetDto struct {
	Id           string     `json:"type" validate:"required"`
	Type         string     `json:"type" validate:"required"`
	Species      string     `json:"species" validate:"required"`
	Name         string     `json:"name" validate:"required"`
	Birthdate    string     `json:"birthdate" validate:"required"`
	Gender       pet.Gender `json:"gender" validate:"required" example:"male"`
	Habit        string     `json:"habit" validate:"required"`
	Caption      string     `json:"caption"`
	Status       pet.Status `json:"status" validate:"required" example:"findhome"`
	IsSterile    bool       `json:"is_sterile" validate:"required"`
	IsVaccinated bool       `json:"is_vaccine" validate:"required"`
	IsVisible    bool       `json:"is_visible" validate:"required"`
	IsClubPet    bool       `json:"is_club_pet" validate:"required"`
	Background   string     `json:"background"`
	Address      string     `json:"address"`
	Contact      string     `json:"contact"`
}

type CreatePetDto struct {
	Pet *PetDto `json:"pet" validate:"required"`
}

type FindOnePetDto struct {
	Id string `json:"id" validate:"required"`
}

type ChangeViewPetReqDto struct {
	Id      string `json:"id" validate:"required"`
	Visible bool   `json:"visible" validate:"required"`
}

type UpDatePetDto struct {
	Type         string     `json:"type" validate:"required"`
	Species      string     `json:"species" validate:"required"`
	Name         string     `json:"name" validate:"required"`
	Birthdate    string     `json:"birthdate" validate:"required"`
	Gender       pet.Gender `json:"gender" validate:"required" example:"male"`
	Habit        string     `json:"habit" validate:"required"`
	Caption      string     `json:"caption"`
	Status       pet.Status `json:"status" validate:"required" example:"findhome"`
	IsSterile    bool       `json:"is_sterile" validate:"required"`
	IsVaccinated bool       `json:"is_vaccine" validate:"required"`
	IsVisible    bool       `json:"is_visible" validate:"required"`
	IsClubPet    bool       `json:"is_club_pet" validate:"required"`
	Background   string     `json:"background"`
	Address      string     `json:"address"`
	Contact      string     `json:"contact"`
}

type DeleteRequest struct {
	Id string `json:"id" validate:"required"`
}
