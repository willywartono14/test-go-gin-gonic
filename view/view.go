package view

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/willywartono14/test-go-gin-gonic/controller"
	"github.com/willywartono14/test-go-gin-gonic/model"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type customerRequest struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type orderRequest struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type OrderRequestTransaction struct {
	Status        string               `json:"status"`
	InvoiceNumber string               `json:"invoice_number"`
	Item          []detailOrderRequest `json:"transaction"`
}

type detailOrderRequest struct {
	ItemName     string `json:"item_name"`
	ItemPrice    int64  `json:"item_price"`
	ItemQuantity int    `json:"item_quantity"`
}

type customerDetailResponse struct {
	Username string                        `json:"username"`
	Fullname string                        `json:"fullname"`
	Email    string                        `json:"email"`
	Phone    string                        `json:"phone"`
	Orders   []customerOrderDetailResponse `json:"orders"`
}

type customerOrderDetailResponse struct {
	OrderID       int    `json:"order_id"`
	Status        string `json:"status"`
	InvoiceNumber string `json:"invoice_number"`
}

type orderDetailResponse struct {
	OrderID       int                              `json:"order_id"`
	Status        string                           `json:"status"`
	InvoiceNumber string                           `json:"invoice_number"`
	Transaction   []orderTransactionDetailResponse `json:"transaction"`
}

type orderTransactionDetailResponse struct {
	ItemName     string `json:"item_name"`
	ItemPrice    int64  `json:"item_price"`
	ItemQuantity int    `json:"item_quantity"`
}

type view struct {
	controller controller.Controller
}

func NewView(controller controller.Controller, router *gin.Engine) {
	view := view{
		controller: controller,
	}

	//Auth
	router.POST("/login", view.login)
	router.POST("/register", view.register)

	//Customer Management
	router.GET("/customers", view.getDataCustomers)
	router.PUT("/customers", view.updateCustomer)
	router.DELETE("/customers/:id", view.deleteCustomer)
	router.GET("/customers/:id", view.getDataDetailCustomer)

	//Order Management
	router.GET("/orders", view.getDataOrders)
	router.PUT("/orders", view.updateOrder)
	router.DELETE("/orders/:id", view.deleteOrder)
	router.POST("/orders", view.insertOrder)
	router.GET("/orders/:id", view.getDataDetailOrder)

	//Items
	router.GET("/items", view.getAllItem)

}

func (a *view) login(context *gin.Context) {

	var loginRequest LoginRequest

	if err := context.BindJSON(&loginRequest); err != nil {
		return
	}

	token, err := a.controller.Login(context, loginRequest.Username, loginRequest.Password)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to login",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"token":      token.AccessToken,
		"expired_at": token.ExpiredTime.Format(time.RFC3339),
	})
}

func (a *view) register(context *gin.Context) {

	var registerRequest registerRequest

	if err := context.BindJSON(&registerRequest); err != nil {
		return
	}

	token, err := a.controller.Register(context, model.User{
		Username: registerRequest.Username,
		Password: registerRequest.Password,
		Fullname: registerRequest.Fullname,
		Email:    registerRequest.Email,
		Phone:    registerRequest.Phone,
	})
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to register",
			"err":     err,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"token":      token.AccessToken,
		"expired_at": token.ExpiredTime.Format(time.RFC3339),
	})
}

func (a *view) getDataCustomers(context *gin.Context) {

	page, _ := context.GetQuery("page")
	page_size, _ := context.GetQuery("page_size")
	search, _ := context.GetQuery("search")

	token := context.Request.Header["Authorization"]

	pageI, err := strconv.Atoi(page)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
		})
		return
	}
	pageSizeI, err := strconv.Atoi(page_size)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
		})
		return
	}

	users, err := a.controller.GetDataCustomers(context, token[0], pageI, pageSizeI, search)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func (a *view) deleteCustomer(context *gin.Context) {

	id := context.Param("id")

	idI, err := strconv.Atoi(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
		})
		return
	}

	token := context.Request.Header["Authorization"]

	err = a.controller.DeleteCustomer(context, token[0], idI)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"status": "data has been successfully deleted",
	})
}

func (a *view) updateCustomer(context *gin.Context) {

	var customer customerRequest

	if err := context.BindJSON(&customer); err != nil {
		return
	}

	token := context.Request.Header["Authorization"]

	err := a.controller.UpdateCustomer(context, token[0], model.User{
		Fullname: customer.Fullname,
		Email:    customer.Email,
		Phone:    customer.Phone,
	})
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to register",
			"err":     err,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "data has been successfully updated",
	})
}

func (a *view) getDataDetailCustomer(context *gin.Context) {

	id := context.Param("id")
	token := context.Request.Header["Authorization"]

	var response customerDetailResponse

	idI, err := strconv.Atoi(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
			"err":     err,
		})
		return
	}

	customer, err := a.controller.GetDetailCustomer(context, token[0], idI)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
			"err":     err,
		})
		return
	}

	if len(customer) > 0 {
		response.Username = customer[0].Username
		response.Fullname = customer[0].Fullname
		response.Email = customer[0].Email
		response.Phone = customer[0].Phone
		for x := range customer {
			var order customerOrderDetailResponse
			order.OrderID = customer[x].OrderID
			order.Status = customer[x].Status
			order.InvoiceNumber = customer[x].InvoiceNumber

			response.Orders = append(response.Orders, order)
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (a *view) getDataOrders(context *gin.Context) {

	page, _ := context.GetQuery("page")
	page_size, _ := context.GetQuery("page_size")
	search, _ := context.GetQuery("search")

	token := context.Request.Header["Authorization"]

	pageI, err := strconv.Atoi(page)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
			"err":     err,
		})
		return
	}
	pageSizeI, err := strconv.Atoi(page_size)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
			"err":     err,
		})
		return
	}

	orders, err := a.controller.GetDataOrders(context, token[0], pageI, pageSizeI, search)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
			"err":     err,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"data": orders,
	})
}

func (a *view) deleteOrder(context *gin.Context) {

	id := context.Param("id")

	idI, err := strconv.Atoi(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
			"err":     err,
		})
		return
	}

	token := context.Request.Header["Authorization"]

	err = a.controller.DeleteOrder(context, token[0], idI)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
			"err":     err,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"status": "data has been successfully deleted",
	})
}

func (a *view) updateOrder(context *gin.Context) {

	var order orderRequest

	if err := context.BindJSON(&order); err != nil {
		return
	}

	token := context.Request.Header["Authorization"]

	err := a.controller.UpdateOrder(context, token[0], model.Order{
		ID:     order.ID,
		Status: order.Status,
	})
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to register",
			"err":     err,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "data has been successfully updated",
	})
}

func (a *view) insertOrder(context *gin.Context) {

	var order OrderRequestTransaction
	var items []model.RequestDetailOrder

	if err := context.BindJSON(&order); err != nil {
		return
	}

	token := context.Request.Header["Authorization"]

	for x := range order.Item {
		var item model.RequestDetailOrder

		item.ItemName = order.Item[x].ItemName
		item.ItemPrice = order.Item[x].ItemPrice
		item.ItemQuantity = order.Item[x].ItemQuantity

		items = append(items, item)

	}

	err := a.controller.InsertOrder(context, token[0], model.RequestOrder{
		Status:        order.Status,
		InvoiceNumber: order.InvoiceNumber,
		Item:          items,
	})

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to insert data",
			"err":     err,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": "data has been successfully inserted",
	})
}

func (a *view) getDataDetailOrder(context *gin.Context) {

	id := context.Param("id")
	token := context.Request.Header["Authorization"]

	var response orderDetailResponse

	idI, err := strconv.Atoi(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
			"error":   err,
		})
		return
	}

	order, err := a.controller.GetDetailOrder(context, token[0], idI)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
			"error":   err,
		})
		return
	}

	if len(order) > 0 {
		response.OrderID = order[0].ID
		response.Status = order[0].Status
		response.InvoiceNumber = order[0].InvoiceNumber
		for x := range order {
			var transaction orderTransactionDetailResponse
			transaction.ItemName = order[x].ItemName
			transaction.ItemPrice = order[x].ItemPrice
			transaction.ItemQuantity = order[x].ItemQuantity

			response.Transaction = append(response.Transaction, transaction)
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (a *view) getAllItem(context *gin.Context) {

	token := context.Request.Header["Authorization"]

	items, err := a.controller.GetAllItem(context, token[0])
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
			"err":     err,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"data": items,
	})
}
