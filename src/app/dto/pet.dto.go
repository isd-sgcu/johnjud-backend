package dto

import (
	"github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
)

type ImageResponse struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type PetResponse struct {
	Id           string          `json:"id"`
	Type         string          `json:"type"`
	Species      string          `json:"species"`
	Name         string          `json:"name"`
	Birthdate    string          `json:"birthdate"`
	Gender       pet.Gender      `json:"gender"`
	Color        string          `json:"color"`
	Pattern      string          `json:"pattern"`
	Habit        string          `json:"habit"`
	Caption      string          `json:"caption"`
	Status       pet.Status      `json:"status"`
	IsSterile    *bool           `json:"is_sterile"`
	IsVaccinated *bool           `json:"is_vaccinated"`
	IsVisible    *bool           `json:"is_visible"`
	IsClubPet    *bool           `json:"is_club_pet"`
	Origin       string          `json:"origin"`
	Address      string          `json:"address"`
	Contact      string          `json:"contact"`
	AdoptBy      string          `json:"adopt_by"`
	Images       []ImageResponse `json:"images"`
}

type CreatePetRequest struct {
	Type         string     `json:"type" validate:"required"`
	Species      string     `json:"species" validate:"required"`
	Name         string     `json:"name" validate:"required"`
	Birthdate    string     `json:"birthdate" validate:"required"`
	Gender       pet.Gender `json:"gender" validate:"required" example:"male"`
	Color        string     `json:"color" validate:"required"`
	Pattern      string     `json:"pattern" validate:"required"`
	Habit        string     `json:"habit" validate:"required"`
	Caption      string     `json:"caption"`
	Status       pet.Status `json:"status" validate:"required" example:"findhome"`
	IsSterile    *bool      `json:"is_sterile" validate:"required"`
	IsVaccinated *bool      `json:"is_vaccinated" validate:"required"`
	IsVisible    *bool      `json:"is_visible" validate:"required"`
	IsClubPet    *bool      `json:"is_club_pet" validate:"required"`
	Origin       string     `json:"origin" validate:"required"`
	Address      string     `json:"address"`
	Contact      string     `json:"contact"`
	AdoptBy      string     `json:"adopt_by"`
	Images       []string   `json:"images"`
}

type ChangeViewPetRequest struct {
	Visible bool `json:"visible" validate:"required"`
}

type ChangeViewPetResponse struct {
	Success bool `json:"success"`
}

type AdoptByRequest struct {
	UserID string `json:"user_id" validate:"required"`
	PetID  string `json:"pet_id" validate:"required"`
}

type AdoptByResponse struct {
	Success bool `json:"success"`
}

type UpdatePetRequest struct {
	Type         string     `json:"type"`
	Species      string     `json:"species"`
	Name         string     `json:"name"`
	Birthdate    string     `json:"birthdate"`
	Gender       pet.Gender `json:"gender"`
	Color        string     `json:"color"`
	Pattern      string     `json:"pattern"`
	Habit        string     `json:"habit"`
	Caption      string     `json:"caption"`
	Status       pet.Status `json:"status"`
	IsSterile    *bool      `json:"is_sterile"`
	IsVaccinated *bool      `json:"is_vaccinated"`
	IsVisible    *bool      `json:"is_visible"`
	IsClubPet    *bool      `json:"is_club_pet"`
	Origin       string     `json:"origin"`
	Address      string     `json:"address"`
	Contact      string     `json:"contact"`
	AdoptBy      string     `json:"adopt_by"`
	Images       []string   `json:"images"`
}
type DeleteRequest struct {
	Id string `json:"id" validate:"required"`
}
type DeleteResponse struct {
	Success bool `json:"success"`
}
