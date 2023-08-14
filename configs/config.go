package configs

import "PasteBay/pkg/utils/environment"

type Config struct {
	Environment    environment.Config
	ConfigFilePath string
}

func InitConfig() *Config {
	configFilePath := "configs/config.json"
	return &Config{
		Environment:    environment.InitEnvironmentConfig(configFilePath),
		ConfigFilePath: configFilePath,
	}
}
