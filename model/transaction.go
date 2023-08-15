package model

import (
	"context"
	"database/sql"
	"time"
)

type Transaction struct {
	ID           int
	OrderID      int
	ItemName     string
	ItemPrice    int
	ItemQuantity int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

func (t Transaction) InsertTransaction(ctx context.Context, db *sql.DB, transaction Transaction) (int, error) {
	query := `INSERT INTO transactions(
		order_id,
  		item_name,
  		item_price,
		item_quantity
	)VALUES (
		$1,
		$2,
		$3,
		$4
	) RETURNING id`
	err := db.QueryRowContext(
		ctx,
		query,
		transaction.OrderID,
		transaction.ItemName,
		transaction.ItemPrice,
		transaction.ItemQuantity,
	).Scan(&transaction.ID)

	if err != nil {
		return 0, err
	}

	return transaction.ID, nil
}
