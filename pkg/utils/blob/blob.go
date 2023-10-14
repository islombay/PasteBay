package blob

import (
	"PasteBay/pkg/utils/logger/sl"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type BlobStorage struct {
	Path          string
	DirectoryPath string
	Log           *slog.Logger
}

func NewBlobStorage(path string, log *slog.Logger) *BlobStorage {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Error("could not create directory for blob storage", sl.Err(err))
			os.Exit(1)
		}
	}

	return &BlobStorage{
		Path: path,
		Log:  log,
	}
}

func (blob *BlobStorage) Save(content string) (string, error) {
	fileName := strconv.FormatInt(time.Now().Unix(), 10)
	filePath := blob.Path + "/" + fileName + ".txt"

	os.Getwd()
	err := os.WriteFile(filePath, []byte(content), 0755)
	if err != nil {
		blob.Log.Error("Could not create blob file", sl.Err(err))
		return "", err
	}
	return filePath, nil
}
