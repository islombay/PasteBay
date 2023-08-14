package repository

import (
	"PasteBay/pkg/models"
	"gorm.io/gorm"
)

type User interface {
	GetUser(id uint) (models.UserModel, error)
}

type Paste interface {
	GetPaste(id uint) (models.PasteModel, error)
	AddPaste(model models.PasteModel) (uint, error)
	DeletePaste(id uint) error
}

type Repository struct {
	User
	Paste
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Paste: NewPasteRepository(db),
		User:  NewUserRepository(db),
	}
}
