package handlers

import (
	"PasteBay/configs"
	_ "PasteBay/docs"
	"PasteBay/pkg/database"
	auth2 "PasteBay/pkg/httpserver/handlers/auth"
	"PasteBay/pkg/httpserver/handlers/pastes"
	"PasteBay/pkg/httpserver/middlewares"
	"PasteBay/pkg/utils/blob"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	router.Use(middlewares.MiddlewareAuth(r.AuthCfg))

	router.GET("/paste/:alias", pastes.GetPaste(r.Log, r.DB, r.Blob))
	router.POST("/paste/:alias", pastes.GetPaste(r.Log, r.DB, r.Blob))

	api := router.Group("/api")
	{
		// CREATE POST
		api.POST("/paste", pastes.AddPaste(r.Log, r.DB, r.Blob, r.Server.AliasPath))

		// DELETE POST FOR AUTHENTICATED USERS
		api.DELETE("/paste", pastes.DeletePasteHandler(r.Log, r.DB, r.Blob))
	}

	auth := router.Group("/auth")
	{
		auth.POST("/register", auth2.RegisterHandler(r.Log, r.DB))
		auth.POST("/login", auth2.LoginHandler(r.Log, r.DB, r.AuthCfg))
	}

	router.GET("/swagger", swaggerRedirect)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
