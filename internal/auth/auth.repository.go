package auth

import (
	"github.com/isd-sgcu/johnjud-backend/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	Create(auth *model.AuthSession) error
	Delete(id string) error
}

type repositoryImpl struct {
	Db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{Db: db}
}

func (r *repositoryImpl) Create(auth *model.AuthSession) error {
	return r.Db.Create(auth).Error
}

func (r *repositoryImpl) Delete(id string) error {
	return r.Db.Delete(&model.AuthSession{}, "id = ?", id).Error
}
