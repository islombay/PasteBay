package main

import (
	"PasteBay"
	"PasteBay/configs"
	"PasteBay/pkg/database"
	"PasteBay/pkg/httpserver/handlers"
	"PasteBay/pkg/utils/blob"
	"PasteBay/pkg/utils/logger/sl"
	"context"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// @title Paste Bay
// @version 1.0
// @description Simple server for pasting and sharing url link to them

// @host localhost:9874
// @BasePath /

func main() {
	config := configs.InitConfig()

	log := setupLogger(config.Env.Environment)

	if err := godotenv.Load(); err != nil {
		log.Error("error loading env variables:", sl.Err(err))
		os.Exit(1)
	}

	db := database.InitDatabase(database.DatabaseLoad{
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		DBName:   config.DB.DBName,
		SSLMode:  config.DB.SSLMode,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),

		Log: log,
	})

	authCfg := config.Auth
	value, err := strconv.Atoi(os.Getenv("AUTH_TOKENTTL"))
	if err != nil {
		log.Error("could not get tokenttl", sl.Err(err))
		os.Exit(1)
	}
	authCfg.TokenTTL = value
	authCfg.SecretKey = os.Getenv("AUTH_SECRETKEY")

	blobStorage := blob.NewBlobStorage(config.Blob.Path, log)
	engine := handlers.InitRoutes(
		handlers.RouteInit{log, db, blobStorage, config.Server, authCfg},
	)

	srv := new(PasteBay.Server)
	go func() {
		if err := srv.Run(config.Server.Port, engine); err != nil {
			log.Error("failed to start server", sl.Err(err))
			os.Exit(1)
		}
	}()

	log.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Server stopped")

	if err := db.Close(); err != nil {
		log.Error("error occured on database connection stopping:", sl.Err(err))
	}
	log.Info("Database connection closed")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error("error occured on server shutting down", sl.Err(err))
	}
	log.Info("Server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case configs.EnvDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case configs.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}
