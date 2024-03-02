// Package usecase provides interfaces for defining order interactor.
package usecase

import (
	"L0/internal/entity"
	"context"
)

//go:generate mockgen -source=./interfaces.go -destination=usecases_mock.go -package=usecase

// OrderInteractor defines the interface for order use cases.
type OrderInteractor interface {
	// Create creates a new order.
	// It takes a context and an order entity as input parameters.
	// Returns an error if the operation fails.
	Create(ctx context.Context, order *entity.Order) error

	// GetByUid retrieves an order by its UID.
	// It takes a context and a UID string as input parameters.
	// Returns the order entity or an error if the operation fails.
	GetByUid(ctx context.Context, uid string) (*entity.Order, error)

	// GetAll retrieves all orders.
	// It takes a context as an input parameter.
	// Returns a slice of order entities or an error if the operation fails.
	GetAll(ctx context.Context) ([]*entity.Order, error)

	// Delete deletes an order.
	// It takes a context and a UID string as input parameters.
	// Returns an error if the operation fails.
	Delete(ctx context.Context, uid string) error
}
