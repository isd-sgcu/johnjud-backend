package dto

import (
	"github.com/isd-sgcu/johnjud-backend/constant"
)

type PetResponse struct {
	Id           string           `json:"id"`
	Type         string           `json:"type"`
	Name         string           `json:"name"`
	Birthdate    string           `json:"birthdate"`
	Gender       constant.Gender  `json:"gender"`
	Color        string           `json:"color"`
	Pattern      string           `json:"pattern"`
	Habit        string           `json:"habit"`
	Caption      string           `json:"caption"`
	Status       constant.Status  `json:"status"`
	IsSterile    *bool            `json:"is_sterile"`
	IsVaccinated *bool            `json:"is_vaccinated"`
	IsVisible    *bool            `json:"is_visible"`
	Origin       string           `json:"origin"`
	Owner        string           `json:"owner"`
	Contact      string           `json:"contact"`
	Tel          string           `json:"tel"`
	Images       []*ImageResponse `json:"images"`
}

type FindAllPetRequest struct {
	Search   string `json:"search"`
	Type     string `json:"type"`
	Gender   string `json:"gender"`
	Color    string `json:"color"`
	Pattern  string `json:"pattern"`
	Age      string `json:"age"`
	MinAge   int    `json:"min_age"`
	MaxAge   int    `json:"max_age"`
	Origin   string `json:"origin"`
	PageSize int    `json:"page_size"`
	Page     int    `json:"page"`
}

type FindAllMetadata struct {
	Page       int `json:"page"`
	TotalPages int `json:"total_pages"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
}

type FindAllPetResponse struct {
	Pets     []*PetResponse   `json:"pets"`
	Metadata *FindAllMetadata `json:"metadata"`
}

type CreatePetRequest struct {
	Type         string          `json:"type" validate:"required"`
	Name         string          `json:"name" validate:"required"`
	Birthdate    string          `json:"birthdate" validate:"required"`
	Gender       constant.Gender `json:"gender" validate:"required" example:"male"`
	Color        string          `json:"color" validate:"required"`
	Pattern      string          `json:"pattern" validate:"required"`
	Habit        string          `json:"habit" validate:"required"`
	Caption      string          `json:"caption"`
	Status       constant.Status `json:"status" validate:"required" example:"findhome"`
	IsSterile    *bool           `json:"is_sterile" validate:"required"`
	IsVaccinated *bool           `json:"is_vaccinated" validate:"required"`
	IsVisible    *bool           `json:"is_visible" validate:"required"`
	Origin       string          `json:"origin" validate:"required"`
	Owner        string          `json:"owner"`
	Contact      string          `json:"contact"`
	Tel          string          `json:"tel"`
	Images       []string        `json:"images"`
}

type ChangeViewPetRequest struct {
	Visible bool `json:"visible"`
}

type ChangeViewPetResponse struct {
	Success bool `json:"success"`
}

type AdoptByRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

type AdoptByResponse struct {
	Success bool `json:"success"`
}

type UpdatePetRequest struct {
	Type         string          `json:"type"`
	Name         string          `json:"name"`
	Birthdate    string          `json:"birthdate"`
	Gender       constant.Gender `json:"gender"`
	Color        string          `json:"color"`
	Pattern      string          `json:"pattern"`
	Habit        string          `json:"habit"`
	Caption      string          `json:"caption"`
	Status       constant.Status `json:"status"`
	IsSterile    *bool           `json:"is_sterile"`
	IsVaccinated *bool           `json:"is_vaccinated"`
	IsVisible    *bool           `json:"is_visible"`
	Origin       string          `json:"origin"`
	Owner        string          `json:"owner"`
	Contact      string          `json:"contact"`
	Tel          string          `json:"tel"`
	Images       []string        `json:"images"`
}
type DeleteRequest struct {
	Id string `json:"id" validate:"required"`
}
type DeleteResponse struct {
	Success bool `json:"success"`
}
