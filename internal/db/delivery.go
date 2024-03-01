// Package db provides methods for working with deliveries in the database.

package db

import (
	"L0/internal/entity"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// CreateDelivery creates a new delivery record in the database.
// It takes a context and a delivery entity as input parameters.
// Returns the unique identifier of the created delivery or an error if the operation fails.
func (s *source) CreateDelivery(ctx context.Context, delivery *entity.Delivery) (string, error) {
	// Create a database context with a timeout
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	// Generate a unique identifier for the delivery
	deliveryUID := uuid.New().String()

	// Execute the SQL query to insert the delivery record into the database
	row := s.db.QueryRowContext(
		dbCtx,
		`INSERT INTO deliveries 
			(delivery_uid, name, phone, zip, city, address, region, email) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)`,
		deliveryUID,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
	)
	if err := row.Err(); err != nil {
		return "", fmt.Errorf("can't execute query: %w", err)
	}

	// Return the unique identifier of the created delivery
	return deliveryUID, nil
}

// GetDeliveryById retrieves a delivery record from the database by its unique identifier.
// It takes a context and a delivery ID as input parameters.
// Returns the delivery record or an error if the operation fails.
func (s *source) GetDeliveryById(ctx context.Context, id string) (*entity.DeliveryDB, error) {
	// Create a database context with a timeout
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	// Execute the SQL query to retrieve the delivery record by ID from the database
	row := s.db.QueryRowxContext(
		dbCtx,
		"SELECT * FROM deliveries WHERE delivery_uid = $1",
		id,
	)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("can't execute query: %w", err)
	}

	// Scan the delivery record from the database into a struct
	var delivery entity.DeliveryDB
	if err := row.StructScan(&delivery); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("can't scan delivery: %w", err)
	}

	// Return the delivery record
	return &delivery, nil
}

// GetDeliveryByEmail retrieves a delivery record from the database by recipient's email address.
// It takes a context and an email address as input parameters.
// Returns the delivery record or an error if the operation fails.
func (s *source) GetDeliveryByEmail(ctx context.Context, email string) (*entity.DeliveryDB, error) {
	// Create a database context with a timeout
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	// Execute the SQL query to retrieve the delivery record by email from the database
	row := s.db.QueryRowxContext(
		dbCtx,
		"SELECT * FROM deliveries WHERE email = $1",
		email,
	)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("can't execute query: %w", err)
	}

	// Scan the delivery record from the database into a struct
	var delivery entity.DeliveryDB
	if err := row.StructScan(&delivery); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("can't scan delivery: %w", err)
	}

	// Return the delivery record
	return &delivery, nil
}

// GetDeliveryByPhone retrieves a delivery record from the database by recipient's phone number.
// It takes a context and a phone number as input parameters.
// Returns the delivery record or an error if the operation fails.
func (s *source) GetDeliveryByPhone(ctx context.Context, phone string) (*entity.DeliveryDB, error) {
	// Create a database context with a timeout
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	// Execute the SQL query to retrieve the delivery record by phone number from the database
	row := s.db.QueryRowxContext(
		dbCtx,
		"SELECT * FROM deliveries WHERE phone = $1",
		phone,
	)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("can't execute query: %w", err)
	}

	// Scan the delivery record from the database into a struct
	var delivery entity.DeliveryDB
	if err := row.StructScan(&delivery); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("can't scan delivery: %w", err)
	}

	// Return the delivery record
	return &delivery, nil
}
