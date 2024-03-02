package handlers

import (
	"L0/internal/nats"
	"L0/internal/usecase"
	"L0/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// orderHandlers represents the implementation of OrderHandlers interface.
type orderHandlers struct {
	interactor  usecase.OrderInteractor
	natsService nats.NATSService
}

// NewOrderHandlers creates a new instance of orderHandlers.
func NewOrderHandlers(interactor usecase.OrderInteractor, natsService nats.NATSService) *orderHandlers {
	return &orderHandlers{
		interactor:  interactor,
		natsService: natsService,
	}
}

// CreateOrderHandler handles requests to create an order.
func (h *orderHandlers) CreateHandler(c *gin.Context) {
	order := utils.GenerateOrder()

	// Marshal the order into JSON format
	orderJSON, err := json.Marshal(order)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't marshal order: %w", err))
		return
	}

	// Publish the order to NATS
	err = h.natsService.Publish(orderJSON)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't publish order: %w", err))
	}
	c.JSON(http.StatusOK, order)
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

// DeleteHandler handles requests to delete an order.
func (h *orderHandlers) DeleteHandler(c *gin.Context) {
	ctx := context.Background()

	uid := c.Param("uid")

	err := h.interactor.Delete(ctx, uid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't delete order: %w", err))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllHandler handles requests to retrieve all orders.
func (h *orderHandlers) GetAllHandler(c *gin.Context) {
	ctx := context.Background()

	orders, err := h.interactor.GetAll(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't get orders: %w", err))
		return
	}

	var ids []string
	// ids := make([]string, len(orders))
	for _, order := range orders {
		ids = append(ids, order.OrderUID)
	}
	fmt.Println(ids)

	c.JSON(http.StatusOK, ids)
}

// GetHTMLOrderHandler handles requests to retrieve an HTML representation of an order.
func (h *orderHandlers) GetHTMLOrderHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "order.html", nil)
}
