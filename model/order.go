package model

import (
	"context"
	"database/sql"
	"time"
)

type Order struct {
	ID            int
	UserID        int
	Status        string
	InvoiceNumber string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}

func (o Order) FindOrdersWithPagination(ctx context.Context, db *sql.DB, page, pageSize int, search string) ([]Order, error) {
	query := `SELECT 
		id, 
		user_id,
		invoice_number,
		status
		FROM orders
		where deleted_at is null`

	if search != "" {
		query += ` and invoice_number ilike '%%'||$3||'%%'`
	}

	query += ` limit $1
			offset $2`

	var results []Order

	var rows *sql.Rows
	var err error

	if search != "" {
		rows, err = db.QueryContext(ctx, query, pageSize, page, search)
	} else {
		rows, err = db.QueryContext(ctx, query, pageSize, page)
	}

	if err != nil {
		return results, err
	}

	for rows.Next() {
		var result Order

		err := rows.Scan(
			&result.ID,
			&result.UserID,
			&result.InvoiceNumber,
			&result.Status,
		)
		if err != nil {
			return results, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (o Order) UpdateOrder(ctx context.Context, db *sql.DB, order Order) error {
	query := `UPDATE orders
		SET status = $2,
			updated_at = $3
	WHERE id = $1`

	_, err := db.ExecContext(
		ctx,
		query,
		order.ID,
		order.Status,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (o Order) DeleteOrder(ctx context.Context, db *sql.DB, id int) error {
	query := `UPDATE orders
		SET deleted_at = $2
	WHERE id = $1`

	_, err := db.ExecContext(
		ctx,
		query,
		id,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (o Order) FindOrderById(ctx context.Context, db *sql.DB, id int64) (Order, error) {
	query := `SELECT id, user_id, status, invoice_number
	FROM orders where id = $1`
	var order Order
	err := db.QueryRowContext(ctx, query, id).
		Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.InvoiceNumber,
		)
	if err != nil && err != sql.ErrNoRows {
		return Order{}, err
	}

	return order, nil
}

func (o Order) InsertOrder(ctx context.Context, db *sql.DB, order Order) (int, error) {
	query := `INSERT INTO orders(
		user_id,
  		status,
  		invoice_number
	)VALUES (
		$1,
		$2,
		$3
	) RETURNING id`
	err := db.QueryRowContext(
		ctx,
		query,
		order.UserID,
		order.Status,
		order.InvoiceNumber,
	).Scan(&order.ID)

	if err != nil {
		return 0, err
	}

	return order.ID, nil
}
