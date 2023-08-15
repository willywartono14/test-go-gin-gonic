package model

import "time"

type Token struct {
	AccessToken string
	ExpiredTime time.Time
}

type ResponseUser struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type ResponseOrder struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	Status        string `json:"status"`
	InvoiceNumber string `json:"invoice_number"`
}

type RequestOrder struct {
	ID            int                  `json:"id"`
	UserID        int                  `json:"user_id"`
	Status        string               `json:"status"`
	InvoiceNumber string               `json:"invoice_number"`
	Item          []RequestDetailOrder `json:"transaction"`
}

type RequestDetailOrder struct {
	ItemName     string `json:"item_name"`
	ItemPrice    int64  `json:"item_price"`
	ItemQuantity int    `json:"item_quantity"`
}

type ResponseUserDetail struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	Fullname      string `json:"fullname"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	OrderID       int    `json:"order_id"`
	Status        string `json:"status"`
	InvoiceNumber string `json:"invoice_number"`
}

type ResponseOrderDetail struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	Status        string `json:"status"`
	InvoiceNumber string `json:"invoice_number"`
	ItemName      string `json:"item_name"`
	ItemPrice     int64  `json:"item_price"`
	ItemQuantity  int    `json:"item_quantity"`
}

type ResponseItem struct {
	ID        int    `json:"id"`
	ItemName  string `json:"item_name"`
	ItemPrice int64  `json:"item_price"`
	ItemStock int    `json:"item_stock"`
}
