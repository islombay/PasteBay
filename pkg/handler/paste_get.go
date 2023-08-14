package handler

import (
	"PasteBay/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) pasteGet(c *gin.Context) {
	hash := c.Param("hash")
	var reqKey models.RequestGetPaste

	if c.Request.Method == "POST" {
		if err := c.BindJSON(&reqKey); err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	pasteObj, err := h.service.Paste.GetPaste(hash, reqKey)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, newResponse(
		"",
		http.StatusOK,
		pasteObj,
		""))
}
