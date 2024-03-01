// Package db provides methods for working with payments in the database.

package db

import (
	"L0/internal/entity"
	"context"
	"database/sql"
	"fmt"
)

// CreatePayment inserts a new payment record into the database.
// It takes a context and a payment entity as input parameters.
// Returns the transaction ID of the created payment or an error if the operation fails.
func (s *source) CreatePayment(ctx context.Context, payment *entity.Payment) (string, error) {
	// Create a database context with a timeout
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	// Execute the query to insert payment details into the database
	row := s.db.QueryRowContext(
		dbCtx,
		`INSERT INTO payments 
		(transaction, request_id, currency, provider, amount, payment_dt, 
		bank, delivery_cost, goods_total, custom_fee) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDt,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	)

	// Check for any errors after executing the query
	if err := row.Err(); err != nil {
		return "", fmt.Errorf("can't execute query: %w", err)
	}

	// Return the transaction ID of the created payment
	return payment.Transaction, nil
}

// GetPaymentByTransaction retrieves payment details from the database by transaction ID.
// It takes a context and a transaction ID as input parameters.
// Returns the payment entity or an error if the operation fails.
func (s *source) GetPaymentByTransaction(ctx context.Context, transaction string) (*entity.Payment, error) {
	// Create a database context with a timeout
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	// Query payment details from the database by transaction ID
	row := s.db.QueryRowxContext(
		dbCtx,
		"SELECT * FROM payments WHERE transaction = $1",
		transaction,
	)
	// Check for any errors after executing the query
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("can't execute query: %w", err)
	}

	// Initialize a variable to store payment details
	var payment entity.Payment

	// Scan the payment record from the database into a struct
	if err := row.StructScan(&payment); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("can't scan payment: %w", err)
	}

	// Return the payment entity
	return &payment, nil
}
