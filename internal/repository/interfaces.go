// Package repository provides interfaces for interacting with order repository.
package repository

import (
	"L0/internal/entity"
	"context"
)

//go:generate mockgen -source=./interfaces.go -destination=repositories_mock.go -package=repository

// OrderRepository defines the interface for order repositories.
type OrderRepository interface {
	// Create inserts a new order into the repository.
	// It takes a context and an order entity as input parameters.
	// Returns the UID of the created order or an error if the operation fails.
	Create(ctx context.Context, order *entity.Order) (string, error)

	// GetByUid retrieves an order from the repository by its UID.
	// It takes a context and a UID string as input parameters.
	// Returns the order entity or an error if the operation fails.
	GetByUid(ctx context.Context, uid string) (*entity.Order, error)

	// GetAll retrieves all orders from the repository.
	// It takes a context as an input parameter.
	// Returns a slice of order entities or an error if the operation fails.
	GetAll(ctx context.Context) ([]*entity.Order, error)

	// Delete deletes an order.
	// It takes a context and a UID string as input parameters.
	// Returns an error if the operation fails.
	Delete(ctx context.Context, uid string) error
}
