package response

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Message string `json:"message"`
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
