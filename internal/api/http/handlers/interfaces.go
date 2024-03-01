// Package handlers provides interfaces for order handlers.
package handlers

import "github.com/gin-gonic/gin"

// OrderHandlers defines the interface for order handlers.
type OrderHandlers interface {
	// GetByIdHandler handles requests to retrieve an order by its ID.
	GetByIdHandler(c *gin.Context)

	// GetHTMLOrderHandler handles requests to retrieve an HTML representation of an order.
	GetHTMLOrderHandler(c *gin.Context)
}
