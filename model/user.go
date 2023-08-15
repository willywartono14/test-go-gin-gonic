package model

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID        int
	Username  string
	Password  string
	Fullname  string
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (u User) FindUserByUsername(ctx context.Context, db *sql.DB, username string) (User, error) {
	query := `SELECT id, username, password, created_at 
	FROM users where username = $1`
	var user User
	err := db.QueryRowContext(ctx, query, username).
		Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
		)
	if err != nil && err != sql.ErrNoRows {
		return User{}, err
	}

	return user, nil
}

func (u User) InsertUser(ctx context.Context, db *sql.DB, user User) (int, error) {
	query := `INSERT INTO users(
		username,
  		password,
  		fullname,
  		email,
  		phone
	)VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	) RETURNING id`
	err := db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Fullname,
		user.Email,
		user.Phone,
	).Scan(&user.ID)

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (u User) FindCustomersWithPagination(ctx context.Context, db *sql.DB, page, pageSize int, search string) ([]User, error) {
	query := `SELECT 
		id, 
		fullname,
		email,
		phone
		FROM users
		where deleted_at is null`

	if search != "" {
		query += ` and fullname ilike '%%'||$3||'%%'`
	}

	query += ` limit $1
		offset $2`

	var results []User

	var rows *sql.Rows
	var err error

	if search != "" {
		rows, err = db.QueryContext(ctx, query, pageSize, page, search)
	} else {
		rows, err = db.QueryContext(ctx, query, pageSize, page)
	}

	defer func() {
		err = rows.Close()
	}()

	if err != nil {
		return results, err
	}

	for rows.Next() {
		var result User

		err := rows.Scan(
			&result.ID,
			&result.Fullname,
			&result.Email,
			&result.Phone,
		)
		if err != nil {
			return results, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (u User) UpdateCustomer(ctx context.Context, db *sql.DB, user User) error {
	query := `UPDATE users
		SET fullname = $2,
			email = $3,
			phone = $4,
			updated_at = $5
	WHERE id = $1`

	_, err := db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Fullname,
		user.Email,
		user.Phone,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (u User) DeleteCustomer(ctx context.Context, db *sql.DB, id int) error {
	query := `UPDATE users
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

func (u User) FindUserById(ctx context.Context, db *sql.DB, id int64) (User, error) {
	query := `SELECT id, username, password, fullname, email, phone  
	FROM users where id = $1`
	var user User
	err := db.QueryRowContext(ctx, query, id).
		Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.Fullname,
			&user.Email,
			&user.Phone,
		)
	if err != nil && err != sql.ErrNoRows {
		return User{}, err
	}

	return user, nil
}

func (u User) GetDetailCustomer(ctx context.Context, db *sql.DB, id int) ([]ResponseUserDetail, error) {
	query := `SELECT u.username, u.fullname, u.email , u.phone , o.id, o.status, o.invoice_number 
	FROM users u
	JOIN orders o
	ON u.id = o.user_id
	WHERE o.deleted_at IS NULL
	AND u.id = 2`
	var results []ResponseUserDetail

	rows, err := db.QueryContext(ctx, query, id)

	defer func() {
		err = rows.Close()
	}()

	if err != nil {
		return results, err
	}

	for rows.Next() {
		var result ResponseUserDetail

		err := rows.Scan(
			&result.ID,
			&result.Fullname,
			&result.Email,
			&result.Phone,
		)
		if err != nil {
			return results, err
		}

		results = append(results, result)
	}

	return results, nil
}
