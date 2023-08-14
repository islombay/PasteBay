package environment

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type Config struct {
	Environment string `json:"environment"`
}

func InitEnvironmentConfig(configPath string) Config {
	currentDir, err := os.Getwd()
	if err != nil {
		logrus.Fatalf("error getting current directory: %s", err.Error())
	}
	//logrus.Info(currentDir)

	configFilePath := filepath.Join(currentDir, configPath)
	configFile, err := os.Open(configFilePath)
	if err != nil {
		logrus.Fatalf("Error opening config.json file: %s", err.Error())
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		logrus.Fatalf("Error decoding config.json: %s", err.Error())
	}

	return config
}
