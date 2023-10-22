package blob

import (
	"PasteBay/pkg/utils/logger/sl"
	"bufio"
	"fmt"
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

func (blob *BlobStorage) checkBlobPath() {
	if _, err := os.Stat(blob.Path); os.IsNotExist(err) {
		if err := os.MkdirAll(blob.Path, os.ModePerm); err != nil {
			blob.Log.Error("could not create directory for blob storage", sl.Err(err))
			os.Exit(1)
		}
	}
}

func (blob *BlobStorage) Save(content string) (string, error) {
	fileName := strconv.FormatInt(time.Now().Unix(), 10)
	filePath := blob.Path + "/" + fileName + ".txt"
	blob.checkBlobPath()

	os.Getwd()
	err := os.WriteFile(filePath, []byte(content), 0755)
	if err != nil {
		blob.Log.Error("Could not create blob file", sl.Err(err))
		return "", err
	}
	return filePath, nil
}

func (blob *BlobStorage) Delete(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		blob.Log.Error(fmt.Sprintf("[BLOB] - Could not delete file (%s) in blob storage", filePath), sl.Err(err))
		return err
	}
	blob.Log.Debug(fmt.Sprintf("[BLOB] - File (%s) in blob storage was deleted", filePath))
	return nil
}

func (blob *BlobStorage) GetContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	content := ""
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	file.Close()
	return content, nil
}
