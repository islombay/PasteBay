package service

import (
	"PasteBay/pkg/models"
	"PasteBay/pkg/repository"
)

type User interface{}

type Paste interface {
	GetPaste(hash string, passKeyList ...models.RequestGetPaste) (models.ResponsePasteView, error)
	AddPaste(model models.RequestAddPaste) (string, error)
}

type Service struct {
	User
	Paste
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Paste: NewPasteService(repos.Paste, *repos),
	}
}
