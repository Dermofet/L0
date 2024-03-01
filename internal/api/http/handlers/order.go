package handlers

import (
	"L0/internal/usecase"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// orderHandlers represents the implementation of OrderHandlers interface.
type orderHandlers struct {
	interactor usecase.OrderInteractor
}

// NewOrderHandlers creates a new instance of orderHandlers.
func NewOrderHandlers(interactor usecase.OrderInteractor) *orderHandlers {
	return &orderHandlers{
		interactor: interactor,
	}
}

// GetByIdHandler handles requests to retrieve an order by its ID.
func (h *orderHandlers) GetByIdHandler(c *gin.Context) {
	ctx := context.Background()

	uid := c.Param("uid")

	order, err := h.interactor.GetByUid(ctx, uid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't get order: %w", err))
		return
	}

	if order == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetHTMLOrderHandler handles requests to retrieve an HTML representation of an order.
func (h *orderHandlers) GetHTMLOrderHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "order.html", nil)
}
