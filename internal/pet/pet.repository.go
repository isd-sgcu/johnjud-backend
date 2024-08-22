package pet

import (
	"errors"

	"github.com/isd-sgcu/johnjud-gateway/internal/model"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(result *[]*model.Pet, isAdmin bool) error {
	if isAdmin {
		return r.db.Model(&model.Pet{}).Find(result).Error
	}
	return r.db.Model(&model.Pet{}).Find(result, "is_visible = ?", true).Error
}

func (r *Repository) FindOne(id string, result *model.Pet) error {
	return r.db.Model(&model.Pet{}).First(result, "id = ?", id).Error
}

func (r *Repository) Create(in *model.Pet) error {
	return r.db.Create(&in).Error
}

func (r *Repository) Update(id string, result *model.Pet) error {
	updateMap := UpdateMap(result)
	return r.db.Model(&result).Updates(updateMap).First(&result, "id = ?", id).Error
}

func (r *Repository) Delete(id string) error {
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
