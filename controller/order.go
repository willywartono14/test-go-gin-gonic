package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/willywartono14/test-go-gin-gonic/model"
)

type OrderController interface {
	GetDataOrders(ctx *gin.Context, token string, page, pageSize int, search string) ([]model.ResponseOrder, error)
	UpdateOrder(ctx *gin.Context, token string, orderRequest model.Order) error
	DeleteOrder(ctx *gin.Context, token string, id int) error
	InsertOrder(ctx *gin.Context, token string, orderRequest model.RequestOrder) error
	GetDetailOrder(ctx *gin.Context, token string, id int) ([]model.ResponseOrderDetail, error)
}

func (c *controller) GetDataOrders(ctx *gin.Context, token string, page, pageSize int, search string) ([]model.ResponseOrder, error) {

	var responseUser []model.ResponseOrder

	_, err := c.verifyToken(token)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	if page == 0 {
		page = 1
	}
	if pageSize == 0 || pageSize > 20 {
		pageSize = 10
	}

	page = (page - 1) * pageSize

	orders, err := model.Order{}.FindOrdersWithPagination(ctx.Request.Context(), c.db, page, pageSize, search)
	if err != nil {
		return nil, ErrCredentialInvalid
	}

	if len(orders) > 0 {
		for x := range orders {
			responseUser = append(responseUser, model.ResponseOrder{
				ID:            orders[x].ID,
				UserID:        orders[x].UserID,
				Status:        orders[x].Status,
				InvoiceNumber: orders[x].InvoiceNumber,
			})
		}
	}

	return responseUser, nil
}

func (c *controller) UpdateOrder(ctx *gin.Context, token string, orderRequest model.Order) error {

	payload, err := c.verifyToken(token)
	id := payload.UserId
	if err != nil {
		return ErrNotAuthorized
	}

	order, err := model.Order{}.FindOrderById(ctx.Request.Context(), c.db, id)
	if err != nil {
		return ErrCredentialInvalid
	}

	response := checkDataOrder(order, orderRequest)

	err = model.Order{}.UpdateOrder(ctx.Request.Context(), c.db, response)
	if err != nil {
		return ErrCredentialInvalid
	}

	return nil
}

func (c *controller) DeleteOrder(ctx *gin.Context, token string, id int) error {

	_, err := c.verifyToken(token)
	if err != nil {
		return ErrNotAuthorized
	}

	err = model.Order{}.DeleteOrder(ctx.Request.Context(), c.db, id)
	if err != nil {
		return ErrCredentialInvalid
	}

	return nil
}

func checkDataOrder(order model.Order, orderRequest model.Order) model.Order {

	if orderRequest.Status == "" {
		orderRequest.Status = order.Status
	}

	return orderRequest
}

func (c *controller) InsertOrder(ctx *gin.Context, token string, orderRequest model.RequestOrder) error {

	payload, err := c.verifyToken(token)
	id := payload.UserId
	if err != nil {
		return ErrNotAuthorized
	}

	response := model.Order{
		UserID:        int(id),
		Status:        orderRequest.Status,
		InvoiceNumber: orderRequest.InvoiceNumber,
	}

	orderId, err := model.Order{}.InsertOrder(ctx.Request.Context(), c.db, response)
	if err != nil {
		return err
	}

	if len(orderRequest.Item) > 0 {
		for x := range orderRequest.Item {
			_, err := model.Transaction{}.InsertTransaction(ctx.Request.Context(), c.db, model.Transaction{
				OrderID:      orderId,
				ItemName:     orderRequest.Item[x].ItemName,
				ItemPrice:    int(orderRequest.Item[x].ItemPrice),
				ItemQuantity: orderRequest.Item[x].ItemQuantity,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *controller) GetDetailOrder(ctx *gin.Context, token string, id int) ([]model.ResponseOrderDetail, error) {

	var responseOrderDetail []model.ResponseOrderDetail

	_, err := c.verifyToken(token)
	if err != nil {
		return responseOrderDetail, ErrNotAuthorized
	}

	responseOrderDetail, err = model.Order{}.GetDetailOrder(ctx.Request.Context(), c.db, int(id))
	if err != nil {
		return nil, err
	}

	return responseOrderDetail, nil
}
