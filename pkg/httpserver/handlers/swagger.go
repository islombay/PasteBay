package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func swaggerRedirect(c *gin.Context) {
	fmt.Println("working")
	c.Redirect(http.StatusFound, "/swagger/index.html")
}
