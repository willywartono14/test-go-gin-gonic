package controller

import (
	"database/sql"
)

type (
	Controller interface {
		AuthController
		CustomerController
		OrderController
		ItemController
	}
	controller struct {
		db *sql.DB
	}
)

func NewController(db *sql.DB) Controller {

	return &controller{db}
}
