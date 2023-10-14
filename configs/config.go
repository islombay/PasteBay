package configs

import (
	"PasteBay/pkg/utils/environment"
	"github.com/spf13/viper"
	"log"
)

const (
	EnvDev  = "development"
	EnvProd = "production"
)

type Config struct {
	Env            environment.Config
	ConfigFilePath string
	Server         ServerConfig
	DB             dbConfig
	Blob           blobStorage
}

type dbConfig struct {
	Host, Port, DBName, SSLMode string
}

type blobStorage struct {
	Path string
}

type ServerConfig struct {
	Port      string
	AliasPath string
}

func InitConfig() *Config {
	configFilePath := "configs/config.json"
	envConf := environment.InitEnvironmentConfig(configFilePath)
	if err := initConfigYaml(envConf.Environment); err != nil {
		log.Fatalf("Could not init config yaml: %s", err.Error())
	}
	return &Config{
		Env:            envConf,
		ConfigFilePath: configFilePath,
		Server: ServerConfig{
			Port:      viper.GetString("server.port"),
			AliasPath: viper.GetString("server.alias_path"),
		},
		DB: dbConfig{
			Host:    viper.GetString("db.host"),
			Port:    viper.GetString("db.port"),
			DBName:  viper.GetString("db.dbname"),
			SSLMode: viper.GetString("db.sslmode"),
		},
		Blob: blobStorage{
			Path: viper.GetString("blob.path"),
		},
	}
}

func initConfigYaml(environ string) error {
	viper.AddConfigPath("configs")
	if environ == EnvDev {
		viper.SetConfigName("development-config")
	} else if environ == EnvProd {
		viper.SetConfigName("producion-config")
	}
	viper.SetConfigType("yaml")

	return viper.ReadInConfig()
}
