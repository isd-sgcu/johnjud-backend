package pet

import (
	"errors"

	"github.com/isd-sgcu/johnjud-gateway/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(result *[]*model.Pet, isAdmin bool) error
	FindOne(id string, result *model.Pet) error
	Create(in *model.Pet) error
	Update(id string, result *model.Pet) error
	Delete(id string) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) FindAll(result *[]*model.Pet, isAdmin bool) error {
	if isAdmin {
		return r.db.Model(&model.Pet{}).Find(result).Error
	}
	return r.db.Model(&model.Pet{}).Find(result, "is_visible = ?", true).Error
}

func (r *repositoryImpl) FindOne(id string, result *model.Pet) error {
	return r.db.Model(&model.Pet{}).First(result, "id = ?", id).Error
}

func (r *repositoryImpl) Create(in *model.Pet) error {
	return r.db.Create(&in).Error
}

func (r *repositoryImpl) Update(id string, result *model.Pet) error {
	updateMap := UpdateMap(result)
	return r.db.Model(&result).Where("id = ?", id).Updates(updateMap).First(&result, "id = ?", id).Error
}

func (r *repositoryImpl) Delete(id string) error {
	var pet model.Pet
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&pet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return r.db.Delete(&pet).Error
}
