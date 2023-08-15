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

type authView struct {
	controller controller.Controller
}

func NewAuthView(controller controller.Controller, router *gin.Engine) {
	authView := authView{
		controller: controller,
	}

	router.POST("/login", authView.login)
	router.POST("/register", authView.register)
	router.GET("/customers", authView.getDataCustomers)
	router.PUT("/customers", authView.updateCustomer)
	router.DELETE("/customers/:id", authView.deleteCustomer)

	router.GET("/orders", authView.getDataOrders)
	router.PUT("/orders", authView.updateOrder)
	router.DELETE("/orders/:id", authView.deleteOrder)
	router.POST("/orders", authView.insertOrder)

}

func (a *authView) login(context *gin.Context) {

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

func (a *authView) register(context *gin.Context) {

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

func (a *authView) getDataCustomers(context *gin.Context) {

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

func (a *authView) deleteCustomer(context *gin.Context) {

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
		"status": "deleted",
	})
}

func (a *authView) updateCustomer(context *gin.Context) {

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
		"status": "success",
	})
}

func (a *authView) getDataOrders(context *gin.Context) {

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

	orders, err := a.controller.GetDataOrders(context, token[0], pageI, pageSizeI, search)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"data": orders,
	})
}

func (a *authView) deleteOrder(context *gin.Context) {

	id := context.Param("id")

	idI, err := strconv.Atoi(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
		})
		return
	}

	token := context.Request.Header["Authorization"]

	err = a.controller.DeleteOrder(context, token[0], idI)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch data",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"status": "deleted",
	})
}

func (a *authView) updateOrder(context *gin.Context) {

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
		"status": "success",
	})
}

func (a *authView) insertOrder(context *gin.Context) {

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
		"status": "success",
	})
}
