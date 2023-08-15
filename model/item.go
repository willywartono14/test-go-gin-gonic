package model

import (
	"context"
	"database/sql"
	"time"
)

type Item struct {
	ID        int
	ItemName  string
	ItemPrice int64
	ItemStock int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (i Item) GetAllItem(ctx context.Context, db *sql.DB) ([]Item, error) {
	query := `SELECT 
		id, 
		item_name,
		item_price,
		item_stock
		FROM items
		where deleted_at is null`

	var results []Item

	rows, err := db.QueryContext(ctx, query)

	defer func() {
		err = rows.Close()
	}()

	if err != nil {
		return results, err
	}

	for rows.Next() {
		var result Item

		err := rows.Scan(
			&result.ID,
			&result.ItemName,
			&result.ItemPrice,
			&result.ItemStock,
		)
		if err != nil {
			return results, err
		}

		results = append(results, result)
	}

	return results, nil
}
