package handler

import (
	"PasteBay/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) addPaste(c *gin.Context) {
	var reqPaste models.RequestAddPaste
	if err := c.BindJSON(&reqPaste); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	paste_hash, err := h.service.Paste.AddPaste(reqPaste)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, newResponse("", http.StatusOK, paste_hash, ""))
}
