package handler

import (
	"PasteBay/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{service: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		paste := api.Group("/paste")
		{
			paste.GET("/:hash", h.pasteGet)
			paste.POST("/:hash", h.pasteGet)
			paste.POST("/", h.addPaste)
		}
	}
	return router
}
