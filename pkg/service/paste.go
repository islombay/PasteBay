package service

import (
	"PasteBay/pkg/models"
	"PasteBay/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const (
	pasteSalt = "P08fzvs2fbx8Yq0cn1ajdziaofnqiUfDKDVUjdbvHF"
)

type PasteService struct {
	repo  repository.Paste
	repos repository.Repository
}

func NewPasteService(repo repository.Paste, repos repository.Repository) *PasteService {
	return &PasteService{repo: repo, repos: repos}
}

func (s *PasteService) GetPaste(hash string, passKeyList ...models.RequestGetPaste) (models.ResponsePasteView, error) {
	var viewObj models.ResponsePasteView
	intHash, err := strconv.Atoi(hash)
	if err != nil {
		logrus.Error(err.Error())
		return viewObj, errors.New("could not convert hash to int")
	}

	pasteModel, err := s.repo.GetPaste(uint(intHash))
	if err != nil {
		return viewObj, err
	}

	if ok := s.checkValidity(pasteModel); !ok {
		if err := s.repo.DeletePaste(pasteModel.ID); err != nil {
			return models.ResponsePasteView{}, errors.New("record not found")
		}
		return models.ResponsePasteView{}, errors.New("record not found")
	}

	if pasteModel.AccessPassword != "" {
		keyProvided := false
		for _, i := range passKeyList {
			if generatePasteHash(i.Password) == pasteModel.AccessPassword {
				keyProvided = true
			}
		}
		if !keyProvided {
			return models.ResponsePasteView{}, errors.New("password is missing")
		}
	}

	author, err := s.repos.User.GetUser(pasteModel.Author)
	if err != nil {
		return viewObj, nil
	}

	viewObj = models.ResponsePasteView{
		ID:        pasteModel.ID,
		CreatedAt: pasteModel.CreatedAt,
		Author: models.ResponseAuthorPasteView{
			Username: author.Username,
		},
		Title:      pasteModel.Title,
		PasteType:  pasteModel.PasteType,
		Content:    pasteModel.BlobSrc,
		ViewsCount: pasteModel.ViewsCount,
	}
	//fmt.Printf("Service later: %#v\n", viewObj)

	s.checkViewLimit(pasteModel)
	return viewObj, nil
}

func (s *PasteService) AddPaste(model models.RequestAddPaste) (string, error) {
	var modelbased models.PasteModel
	changeStructs(&model, &modelbased)

	if model.ExpireTimeMilliseconds == 0 {
		modelbased.ExpireTime = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	} else {
		modelbased.ExpireTime = time.Now().Add(time.Millisecond * time.Duration(model.ExpireTimeMilliseconds))
	}

	modelbased.BlobSrc = model.Content
	modelbased.ViewsCount = 0

	if modelbased.PasteType == 0 {
		modelbased.PasteType = 1
	}

	if len(modelbased.AccessPassword) != 0 {
		modelbased.AccessPassword = generatePasteHash(modelbased.AccessPassword)
	}
	modelbased.Author = 0
	pasteID, err := s.repo.AddPaste(modelbased)

	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(pasteID)), nil
}

func (s *PasteService) checkValidity(model models.PasteModel) bool {
	isValid := true
	currentTime := time.Now()
	defaultTime := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	if !model.ExpireTime.Equal(defaultTime) {
		if model.ExpireTime.Before(currentTime) || model.ExpireTime.Equal(currentTime) {
			isValid = false
		}
	}

	return isValid
}

func (s *PasteService) checkViewLimit(model models.PasteModel) error {
	if model.ViewsLimit != 0 {
		current := model.ViewsCount + 1
		if current == model.ViewsLimit {
			_ = s.repo.DeletePaste(model.ID)
		} else {
			if err := s.repo.UpdateViewsCount(model.ID, current); err != nil {
				logrus.Error(err.Error())
				return err
			}
		}
	}

	return nil
}

func generatePasteHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(pasteSalt)))
}
