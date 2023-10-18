package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ErrorBadRequest         = "bad_request"
	ErrorServerError        = "server_error"
	ErrorNotFound           = "not_found"
	ErrorInvalidCredentials = "invalid_credentials"
	ErrorUserExists         = "username_already_exists"
	ErrorCouldNotCreateUser = "could_not_create_user"
	ErrorAuthPassword       = "invalid_password_type"
	ErrorAuthUsername       = "invalid_username_type"
)

type responsesStruct struct {
	Message    string
	StatusCode int
}

var responses = map[string]responsesStruct{
	ErrorBadRequest: {
		"Bad Request",
		http.StatusBadRequest,
	},
	ErrorServerError: {
		"Server Error",
		http.StatusInternalServerError,
	},
	ErrorNotFound: {
		"Not Found",
		http.StatusNotFound,
	},
	"default": {
		"Unknown Error",
		http.StatusInternalServerError,
	},
	ErrorInvalidCredentials: {
		"Invalid Credentials",
		http.StatusForbidden,
	},
	ErrorUserExists: {
		"Username is already taken",
		http.StatusNotAcceptable,
	},
	ErrorCouldNotCreateUser: {
		"Could not create user",
		http.StatusInternalServerError,
	},
	ErrorAuthPassword: {
		"Password should consist at least of 8 characters",
		http.StatusBadRequest,
	},
	ErrorAuthUsername: {
		"Username should consist of 5 characters, that are: only english letters and underscores",
		http.StatusBadRequest,
	},
}

type errorResponse struct {
	Message string `json:"message"`
}

//func ErrorResponse(c *gin.Context, statusCode int, message string) {
//	c.AbortWithStatusJSON(statusCode, errorResponse{message})
//}

func ErrorResponse(c *gin.Context, msg string) {
	body, ok := responses[msg]
	if !ok {
		body = responses["default"]

	}
	c.AbortWithStatusJSON(body.StatusCode, errorResponse{body.Message})
}
