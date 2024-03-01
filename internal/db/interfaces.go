// Package db provides interfaces and methods for working with the database.

package db

import (
	"L0/internal/entity"
	"context"
)

// OrderSource provides methods for working with orders in the database.
type OrderSource interface {
	// CreateOrder creates a new order in the database.
	// It returns the unique identifier of the created order or an error if the operation fails.
	CreateOrder(ctx context.Context, order *entity.Order) (string, error)

	// GetOrderByUid returns an order from the database by its unique identifier.
	// It returns the order and an error if the order with the specified identifier is not found.
	GetOrderByUid(ctx context.Context, uid string) (*entity.Order, error)

	// GetAllOrders returns all orders from the database.
	// It returns a list of orders and an error if the operation fails.
	GetAllOrders(ctx context.Context) ([]*entity.Order, error)
}

// DeliverySource provides methods for working with deliveries in the database.
type DeliverySource interface {
	// CreateDelivery creates a new delivery record in the database.
	// It returns the unique identifier of the created delivery or an error if the operation fails.
	CreateDelivery(ctx context.Context, delivery *entity.Delivery) (string, error)

	// GetDeliveryById returns a delivery record from the database by its unique identifier.
	// It returns the delivery record and an error if the record with the specified identifier is not found.
	GetDeliveryById(ctx context.Context, id string) (*entity.DeliveryDB, error)

	// GetDeliveryByPhone returns a delivery record from the database by recipient's phone number.
	// It returns the delivery record and an error if the record with the specified phone number is not found.
	GetDeliveryByPhone(ctx context.Context, phone string) (*entity.DeliveryDB, error)

	// GetDeliveryByEmail returns a delivery record from the database by recipient's email address.
	// It returns the delivery record and an error if the record with the specified email address is not found.
	GetDeliveryByEmail(ctx context.Context, email string) (*entity.DeliveryDB, error)
}

// PaymentSource provides methods for working with payments in the database.
type PaymentSource interface {
	// CreatePayment creates a new payment record in the database.
	// It returns the unique identifier of the created payment or an error if the operation fails.
	CreatePayment(ctx context.Context, payment *entity.Payment) (string, error)

	// GetPaymentByTransaction returns a payment record from the database by its transaction identifier.
	// It returns the payment record and an error if the record with the specified transaction identifier is not found.
	GetPaymentByTransaction(ctx context.Context, transaction string) (*entity.Payment, error)
}

// ItemSource provides methods for working with items in the database.
type ItemSource interface {
	// CreateItem creates a new item in the database.
	// It returns the unique identifier of the created item or an error if the operation fails.
	CreateItem(ctx context.Context, item *entity.Item) (string, error)

	// CreateItems creates multiple new items in the database.
	// It returns a list of unique identifiers of the created items or an error if the operation fails.
	CreateItems(ctx context.Context, items []entity.Item) ([]string, error)

	// GetItemByUid returns an item from the database by its unique identifier.
	// It returns the item and an error if the item with the specified identifier is not found.
	GetItemByUid(ctx context.Context, uid string) (*entity.Item, error)

	// GetItemsByTrackNumber returns a list of items from the database by tracking number.
	// It returns a list of items and an error if items with the specified tracking number are not found.
	GetItemsByTrackNumber(ctx context.Context, trackNumber string) ([]entity.Item, error)
}
