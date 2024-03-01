// Package db provides methods for working with orders in the database.

package db

import (
	"L0/internal/entity"
	"context"
	"database/sql"
	"fmt"
)

// CreateOrder creates a new order record in the database.
// It takes a context and an order entity as input parameters.
// Returns the unique identifier of the created order or an error if the operation fails.
func (s *source) CreateOrder(ctx context.Context, order *entity.Order) (string, error) {
	// Create a database context with a timeout
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	// Initialize a variable to store the delivery UID
	var deliveryUID string

	// Check if delivery information is provided
	if order.Delivery.Phone != "" {
		// Get delivery information by phone number
		delivery, err := s.GetDeliveryByPhone(dbCtx, order.Delivery.Phone)
		if err != nil && err != sql.ErrNoRows {
			return "", fmt.Errorf("can't get delivery by phone: %w", err)
		}
		// If delivery information exists, assign its UID
		if delivery != nil {
			deliveryUID = delivery.DeliveryUID
		}
	}

	// If delivery information by phone is not found, check by email
	if order.Delivery.Email != "" && deliveryUID == "" {
		// Get delivery information by email
		delivery, err := s.GetDeliveryByEmail(dbCtx, order.Delivery.Email)
		if err != nil && err != sql.ErrNoRows {
			return "", fmt.Errorf("can't get delivery by email: %w", err)
		}
		// If delivery information exists, assign its UID
		if delivery != nil {
			deliveryUID = delivery.DeliveryUID
		}
	}

	// If delivery information is not found, create a new delivery record
	if deliveryUID == "" {
		var err error
		deliveryUID, err = s.CreateDelivery(dbCtx, &order.Delivery)
		if err != nil {
			return "", fmt.Errorf("can't create delivery: %w", err)
		}
	}

	// Create payment record
	_, err := s.CreatePayment(dbCtx, &order.Payment)
	if err != nil {
		return "", fmt.Errorf("can't create payment: %w", err)
	}

	// Create item records
	_, err = s.CreateItems(dbCtx, order.Items)
	if err != nil {
		return "", fmt.Errorf("can't create items: %w", err)
	}

	// Insert order record into the database
	row := s.db.QueryRowContext(
		dbCtx,
		`INSERT INTO orders
		(order_uid, track_number, entry, delivery_uid, payment_transaction, locale, 
		internal_signature, customer_id, delivery_service, shardkey, sm_id, 
		date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		deliveryUID,
		order.Payment.Transaction,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err := row.Err(); err != nil {
		return "", fmt.Errorf("can't execute query: %w", err)
	}

	// Return the unique identifier of the created order
	return order.OrderUID, nil
}

// GetOrderByUid retrieves an order record from the database by its unique identifier.
// It takes a context and an order UID as input parameters.
// Returns the order record or an error if the operation fails.
func (s *source) GetOrderByUid(ctx context.Context, orderUID string) (*entity.Order, error) {
	// Create a database context with a timeout
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	// Query the order record from the database by UID
	row := s.db.QueryRowxContext(
		dbCtx,
		`SELECT * FROM orders WHERE order_uid = $1`,
		orderUID,
	)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("can't execute query: %w", err)
	}

	// Initialize a variable to store the order details
	var orderDB entity.OrderDB

	// Scan the order record from the database into a struct
	if err := row.StructScan(&orderDB); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("can't scan row: %w", err)
	}

	// Initialize an order entity
	order := entity.Order{
		OrderUID:          orderDB.OrderUID,
		TrackNumber:       orderDB.TrackNumber,
		Entry:             orderDB.Entry,
		Locale:            orderDB.Locale,
		InternalSignature: orderDB.InternalSignature,
		CustomerID:        orderDB.CustomerID,
		DeliveryService:   orderDB.DeliveryService,
		Shardkey:          orderDB.Shardkey,
		SmID:              orderDB.SmID,
		DateCreated:       orderDB.DateCreated,
		OofShard:          orderDB.OofShard,
	}

	// Get delivery details by delivery UID
	delivery, err := s.GetDeliveryById(dbCtx, orderDB.DeliveryUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("can't get delivery: %w", err)
	}
	order.Delivery = entity.Delivery{
		Name:    delivery.Name,
		Phone:   delivery.Phone,
		Zip:     delivery.Zip,
		City:    delivery.City,
		Address: delivery.Address,
		Region:  delivery.Region,
		Email:   delivery.Email,
	}

	// Get payment details by transaction ID
	payment, err := s.GetPaymentByTransaction(dbCtx, orderDB.PaymentTransaction)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("can't get payment: %w", err)
	}
	order.Payment = *payment

	// Get items associated with the order by track number
	items, err := s.GetItemsByTrackNumber(dbCtx, orderDB.TrackNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("can't get items: %w", err)
	}
	order.Items = items

	// Return the order entity
	return &order, nil
}

// GetAllOrders retrieves all order records from the database.
// It takes a context as an input parameter.
// Returns a slice of order records or an error if the operation fails.
func (s *source) GetAllOrders(ctx context.Context) ([]*entity.Order, error) {
	// Create a database context with a timeout
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	// Query all order records from the database
	rows, err := s.db.QueryxContext(
		dbCtx,
		`SELECT * FROM orders`,
	)
	if err != nil {
		return nil, fmt.Errorf("can't execute query: %w", err)
	}
	defer rows.Close()

	// Initialize a slice to store order entities
	var orders []*entity.Order

	// Iterate over each row and scan order details into an entity
	for rows.Next() {
		var orderDB entity.OrderDB
		if err := rows.StructScan(&orderDB); err != nil {
			return nil, fmt.Errorf("can't scan row: %w", err)
		}

		// Initialize an order entity
		order := entity.Order{
			OrderUID:          orderDB.OrderUID,
			TrackNumber:       orderDB.TrackNumber,
			Entry:             orderDB.Entry,
			Locale:            orderDB.Locale,
			InternalSignature: orderDB.InternalSignature,
			CustomerID:        orderDB.CustomerID,
			DeliveryService:   orderDB.DeliveryService,
			Shardkey:          orderDB.Shardkey,
			SmID:              orderDB.SmID,
			DateCreated:       orderDB.DateCreated,
			OofShard:          orderDB.OofShard,
		}

		// Get delivery details by delivery UID
		delivery, err := s.GetDeliveryById(dbCtx, orderDB.DeliveryUID)
		if err != nil {
			return nil, fmt.Errorf("can't get delivery: %w", err)
		}
		order.Delivery = entity.Delivery{
			Name:    delivery.Name,
			Phone:   delivery.Phone,
			Zip:     delivery.Zip,
			City:    delivery.City,
			Address: delivery.Address,
			Region:  delivery.Region,
			Email:   delivery.Email,
		}

		// Get payment details by transaction ID
		payment, err := s.GetPaymentByTransaction(dbCtx, orderDB.PaymentTransaction)
		if err != nil {
			return nil, fmt.Errorf("can't get payment: %w", err)
		}
		order.Payment = *payment

		// Get items associated with the order by track number
		items, err := s.GetItemsByTrackNumber(dbCtx, orderDB.TrackNumber)
		if err != nil {
			return nil, fmt.Errorf("can't get items: %w", err)
		}
		order.Items = items

		// Append the order entity to the slice
		orders = append(orders, &order)
	}

	// Check for errors after iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	// Return the slice of order entities
	return orders, nil
}
