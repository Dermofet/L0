// Package handlers provides interfaces for order handlers.
package handlers

import "github.com/gin-gonic/gin"

//go:generate mockgen -source=interfaces.go -destination=handlers_mock.go -package=handlers

// OrderHandlers defines the interface for order handlers.
type OrderHandlers interface {
	// CreateHandler handles requests to create an order.
	CreateHandler(c *gin.Context)

	// GetByIdHandler handles requests to retrieve an order by its ID.
	GetByIdHandler(c *gin.Context)

	// GetHTMLOrderHandler handles requests to retrieve an HTML representation of an order.
	GetHTMLOrderHandler(c *gin.Context)

	// GetAllHandler handles requests to retrieve all orders.
	GetAllHandler(c *gin.Context)

	// DeleteHandler handles requests to delete an order.
	DeleteHandler(c *gin.Context)
}
