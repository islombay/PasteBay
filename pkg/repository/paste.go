package repository

import (
	"PasteBay/pkg/models"
	"gorm.io/gorm"
)

type PasteRepository struct {
	db *gorm.DB
}

func NewPasteRepository(db *gorm.DB) *PasteRepository {
	return &PasteRepository{db: db}
}

func (r *PasteRepository) GetPaste(id uint) (models.PasteModel, error) {
	var pasteObj models.PasteModel
	if res := r.db.First(&pasteObj, id); res.Error != nil {
		return pasteObj, res.Error
	}
	//fmt.Printf("Repository: %#v\n", pasteObj)

	return pasteObj, nil
}

func (r *PasteRepository) AddPaste(model models.PasteModel) (uint, error) {
	if res := r.db.Create(&model); res.Error != nil {
		return 0, res.Error
	}

	return model.ID, nil
}

func (r *PasteRepository) DeletePaste(id uint) error {
	if res := r.db.Delete(&models.PasteModel{}, id); res.Error != nil {
		return res.Error
	}
	return nil
}
