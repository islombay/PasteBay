package handlers

import (
	"PasteBay/configs"
	"PasteBay/pkg/database"
	"PasteBay/pkg/utils/blob"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type RouteInit struct {
	Log    *slog.Logger
	DB     *database.Database
	Blob   *blob.BlobStorage
	Server configs.ServerConfig
}

func InitRoutes(r RouteInit) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/paste", addPaste(r.Log, r.DB, r.Blob, r.Server.AliasPath))
	}

	return router
}
