package pet

import (
	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(result *[]*pet.Pet) error {
	return r.db.Model(&pet.Pet{}).Find(result).Error
}

func (r *Repository) FindOne(id string, result *pet.Pet) error {
	return r.db.Model(&pet.Pet{}).First(result, "id = ?", id).Error
}

func (r *Repository) Create(in *pet.Pet) error {
	return r.db.Create(&in).Error
}

func (r *Repository) Update(id string, result *pet.Pet) error {
	return r.db.Where(id, "id = ?", id).Updates(&result).First(&result, "id = ?", id).Error
}

func (r *Repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&pet.Pet{}).Error
}
