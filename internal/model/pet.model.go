package model

import (
	"github.com/isd-sgcu/johnjud-backend/constant"
)

type Pet struct {
	Base
	Type         string          `json:"type" gorm:"tinytext"`
	Name         string          `json:"name" gorm:"tinytext"`
	Birthdate    string          `json:"birthdate" gorm:"tinytext"`
	Gender       constant.Gender `json:"gender" gorm:"tinytext" example:"male"`
	Color        string          `json:"color" gorm:"tinytext"`
	Habit        string          `json:"habit" gorm:"mediumtext"`
	Caption      string          `json:"caption" gorm:"mediumtext"`
	Status       constant.Status `json:"status" gorm:"mediumtext" example:"findhome"`
	IsSterile    bool            `json:"is_sterile"`
	IsVaccinated bool            `json:"is_vaccinated"`
	IsVisible    bool            `json:"is_visible"`
	Origin       string          `json:"origin" gorm:"tinytext"`
	Owner        string          `json:"owner" gorm:"tinytext"`
	Contact      string          `json:"contact" gorm:"tinytext"`
	Tel          string          `json:"tel" gorm:"tinytext"`
}
