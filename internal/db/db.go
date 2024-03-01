// Package db provides methods for working with a database.

package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// QueryTimeout specifies the maximum time allowed for a database query to execute.
const (
	QueryTimeout = 10 * time.Second
)

// source represents the data source for interacting with the database.
type source struct {
	db *sqlx.DB
}

// NewSource creates a new instance of the database source with the provided SQLx database connection.
func NewSource(db *sqlx.DB) *source {
	return &source{
		db: db,
	}
}
