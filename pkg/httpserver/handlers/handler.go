package handlers

import (
	"PasteBay/configs"
	"PasteBay/pkg/database"
	auth2 "PasteBay/pkg/httpserver/handlers/auth"
	"PasteBay/pkg/httpserver/handlers/pastes"
	"PasteBay/pkg/utils/blob"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type RouteInit struct {
	Log     *slog.Logger
	DB      *database.Database
	Blob    *blob.BlobStorage
	Server  configs.ServerConfig
	AuthCfg configs.AuthConfig
}

func InitRoutes(r RouteInit) *gin.Engine {
	router := gin.Default()

	router.GET("/paste/:alias", pastes.GetPaste(r.Log, r.DB, r.Blob))
	router.POST("/paste/:alias", pastes.GetPaste(r.Log, r.DB, r.Blob))

	api := router.Group("/api")
	{
		api.POST("/paste", pastes.AddPaste(r.Log, r.DB, r.Blob, r.Server.AliasPath))
	}

	auth := router.Group("/auth")
	{
		auth.POST("/register", auth2.RegisterHandler(r.Log, r.DB))
		auth.POST("/login", auth2.LoginHandler(r.Log, r.DB, r.AuthCfg))
	}

	return router
}
