// Package handlers provides HTTP handler functions.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotImplementedHandler is a handler function for returning a 405 Method Not Allowed status.
func NotImplementedHandler(c *gin.Context) {
	c.AbortWithStatus(http.StatusMethodNotAllowed)
}
