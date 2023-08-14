package main

import (
	"PasteBay"
	"PasteBay/configs"
	"PasteBay/pkg/handler"
	"PasteBay/pkg/repository"
	"PasteBay/pkg/service"
	"PasteBay/pkg/utils/logger"
	"context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := configs.InitConfig()

	logger.InitLogger()
	if err := initConfigYaml(*config); err != nil {
		logrus.Fatalf("error in initiation config.yml files: %s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.DBConfig{
		Port:     viper.GetString("db.port"),
		Host:     viper.GetString("db.host"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("error creating new PostgresDB object: %s", err.Error())
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	srv := new(PasteBay.Server)
	go func() {
		if err := srv.Run(
			viper.GetString("port"),
			handlers.InitRoutes()); err != nil {
			logrus.Fatalf("failed to start server: %s", err.Error())
		}
	}()

	logrus.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Info("Server stopped")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("error occured on server shutting down: %s", err.Error())
	}
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		logrus.Fatalf("error occured on database connection stopping: %s", err.Error())
	}
}

func initConfigYaml(config configs.Config) error {
	viper.AddConfigPath("configs")
	if config.Environment.Environment == "development" {
		viper.SetConfigName("development-config")
		if err := godotenv.Load(); err != nil {
			logrus.Fatalf("error loading env variables: %s", err.Error())
		}
	} else if config.Environment.Environment == "docker" {
		viper.SetConfigName("producion-config")
	}
	viper.SetConfigType("yaml")

	return viper.ReadInConfig()
}
