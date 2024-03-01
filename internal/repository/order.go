// Package repository provides implementations for interacting with order repositories.
package repository

import (
	"L0/internal/db"
	"L0/internal/entity"
	"context"
	"database/sql"
	"fmt"
)

// orderRepository implements the OrderRepository interface.
type orderRepository struct {
	source db.OrderSource
}

// NewOrderRepository creates a new instance of orderRepository.
func NewOrderRepository(source db.OrderSource) *orderRepository {
	return &orderRepository{
		source: source,
	}
}

// Create inserts a new order into the repository.
func (o *orderRepository) Create(ctx context.Context, order *entity.Order) (string, error) {
	id, err := o.source.CreateOrder(ctx, order)
	if err != nil {
		return "", fmt.Errorf("can't create order in db: %w", err)
	}

	return id, nil
}

// GetByUid retrieves an order from the repository by its UID.
func (o *orderRepository) GetByUid(ctx context.Context, uid string) (*entity.Order, error) {
	order, err := o.source.GetOrderByUid(ctx, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get order by uid from db: %w", err)
	}

	return order, nil
}

// GetAll retrieves all orders from the repository.
func (o *orderRepository) GetAll(ctx context.Context) ([]*entity.Order, error) {
	orders, err := o.source.GetAllOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get all orders from db: %w", err)
	}

	return orders, nil
}
