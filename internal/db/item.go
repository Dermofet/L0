// Package db provides methods for working with items in the database.

package db

import (
	"L0/internal/entity"
	"context"
	"database/sql"
	"fmt"
)

// CreateItem creates a new item record in the database.
// It takes a context and an item entity as input parameters.
// Returns the unique identifier of the created item or an error if the operation fails.
func (s *source) CreateItem(ctx context.Context, item *entity.Item) (string, error) {
	// Create a database context with a timeout
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	// Execute the SQL query to insert the item record into the database
	row := s.db.QueryRowContext(
		dbCtx,
		`INSERT INTO items
		(chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		item.ChrtID,
		item.TrackNumber,
		item.Price,
		item.Rid,
		item.Name,
		item.Sale,
		item.Size,
		item.TotalPrice,
		item.NmID,
		item.Brand,
		item.Status,
	)
	if err := row.Err(); err != nil {
		return "", fmt.Errorf("can't execute query: %w", err)
	}

	// Return the unique identifier of the created item
	return item.Rid, nil
}

// CreateItems creates multiple new item records in the database.
// It takes a context and a slice of item entities as input parameters.
// Returns a slice of unique identifiers of the created items or an error if the operation fails.
func (s *source) CreateItems(ctx context.Context, items []entity.Item) ([]string, error) {
	// Create a database context with a timeout
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	// Begin a transaction
	tx, err := s.db.BeginTx(dbCtx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create a slice to store the unique identifiers of the created items
	rids := make([]string, 0, len(items))

	// Prepare the SQL query for inserting items into the database
	stmt, err := tx.PrepareContext(dbCtx,
		`INSERT INTO items
        (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING rid`,
	)
	if err != nil {
		return nil, fmt.Errorf("can't prepare query: %w", err)
	}
	defer stmt.Close()

	// Iterate over each item and execute the SQL query to insert it into the database
	for _, item := range items {
		var rid string
		err := stmt.QueryRowContext(dbCtx,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		).Scan(&rid)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, sql.ErrNoRows
			}
			return nil, fmt.Errorf("can't execute query: %w", err)
		}
		rids = append(rids, rid)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return the unique identifiers of the created items
	return rids, nil
}

// GetItemByUid retrieves an item record from the database by its unique identifier.
// It takes a context and an item ID as input parameters.
// Returns the item record or an error if the operation fails.
func (s *source) GetItemByUid(ctx context.Context, uid string) (*entity.Item, error) {
	// Create a database context with a timeout
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	// Execute the SQL query to retrieve the item record by ID from the database
	row := s.db.QueryRowxContext(
		dbCtx,
		"SELECT * FROM items WHERE chrt_id = $1",
		uid,
	)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("can't execute query: %w", row.Err())
	}

	// Scan the item record from the database into a struct
	var itemDB entity.Item
	if err := row.StructScan(&itemDB); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("can't scan item: %w", err)
	}

	// Return the item record
	return &itemDB, nil
}

// GetItemsByTrackNumber retrieves item records from the database by tracking number.
// It takes a context and a tracking number as input parameters.
// Returns a slice of item records or an error if the operation fails.
func (s *source) GetItemsByTrackNumber(ctx context.Context, trackNumber string) ([]entity.Item, error) {
	// Create a database context with a timeout
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	// Execute the SQL query to retrieve item records by tracking number from the database
	rows, err := s.db.QueryxContext(
		dbCtx,
		"SELECT * FROM items WHERE track_number = $1",
		trackNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("can't execute query: %w", err)
	}
	defer rows.Close()

	// Create a slice to store pointers to item structs temporarily
	var itemPointers []*entity.Item
	for rows.Next() {
		var item entity.Item
		if err := rows.StructScan(&item); err != nil {
			if err == sql.ErrNoRows {
				return nil, sql.ErrNoRows
			}
			return nil, fmt.Errorf("can't scan item: %w", err)
		}
		itemPointers = append(itemPointers, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning items: %w", err)
	}

	// Convert pointers to item structs into item structs before returning
	items := make([]entity.Item, len(itemPointers))
	for i, ptr := range itemPointers {
		items[i] = *ptr
	}

	return items, nil
}
