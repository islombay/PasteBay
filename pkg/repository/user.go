package repository

import (
	"PasteBay/pkg/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(id uint) (models.UserModel, error) {
	var finding models.UserModel
	if res := r.db.First(&finding, id); res.Error != nil {
		return finding, res.Error
	}

	return finding, nil
}
