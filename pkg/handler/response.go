package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StandardResopnse struct {
	Ok     string      `json:"ok"`
	Status int         `json:"status_code"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, newResponse(
		"false",
		statusCode,
		nil,
		message,
	),
	)

	logrus.Errorf("error occured: %s", message)
}

func newResponse(message string, statusCode int, data interface{}, error_msg string) StandardResopnse {
	if message == "" {
		message = "true"
	}
	if error_msg == "" {
		error_msg = "null"
	}
	return StandardResopnse{
		Ok:     message,
		Status: statusCode,
		Data:   data,
		Error:  error_msg,
	}
}
