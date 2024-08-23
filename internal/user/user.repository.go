package user

import (
	"github.com/isd-sgcu/johnjud-backend/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(user *[]*model.User) error
	FindById(id string, user *model.User) error
	FindByEmail(email string, user *model.User) error
	Create(user *model.User) error
	Update(id string, user *model.User) error
	Delete(id string) error
}

type repositoryImpl struct {
	Db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{Db: db}
}

func (r *repositoryImpl) FindAll(user *[]*model.User) error {
	return r.Db.Find(user).Error
}

func (r *repositoryImpl) FindById(id string, user *model.User) error {
	return r.Db.First(user, "id = ?", id).Error
}

func (r *repositoryImpl) FindByEmail(email string, user *model.User) error {
	return r.Db.First(user, "email = ?", email).Error
}

func (r *repositoryImpl) Create(user *model.User) error {
	return r.Db.Create(user).Error
}

func (r *repositoryImpl) Update(id string, user *model.User) error {
	return r.Db.Where("id = ?", id).Updates(user).First(user, "id = ?", id).Error
}

func (r *repositoryImpl) Delete(id string) error {
	return r.Db.Delete(&model.User{}, "id = ?", id).Error
}
