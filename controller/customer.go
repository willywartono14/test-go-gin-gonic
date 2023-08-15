package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/willywartono14/test-go-gin-gonic/model"
)

type CustomerConstroller interface {
	GetDataCustomers(ctx *gin.Context, token string, page, pageSize int, search string) ([]model.ResponseUser, error)
	UpdateCustomer(ctx *gin.Context, token string, userRequest model.User) error
	DeleteCustomer(ctx *gin.Context, token string, id int) error
}

var ErrNotAuthorized error = errors.New("not authorized")

func (c *controller) GetDataCustomers(ctx *gin.Context, token string, page, pageSize int, search string) ([]model.ResponseUser, error) {

	var responseUser []model.ResponseUser

	_, err := c.verifyToken(token)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	users, err := model.User{}.FindCustomersWithPagination(ctx.Request.Context(), c.db, page, pageSize, search)
	if err != nil {
		return nil, ErrCredentialInvalid
	}

	if len(users) > 0 {
		for x := range users {
			responseUser = append(responseUser, model.ResponseUser{
				ID:       users[x].ID,
				Fullname: users[x].Fullname,
				Email:    users[x].Email,
				Phone:    users[x].Phone,
			})
		}
	}

	return responseUser, nil
}

func (c *controller) UpdateCustomer(ctx *gin.Context, token string, userRequest model.User) error {

	payload, err := c.verifyToken(token)
	id := payload.UserId
	if err != nil {
		return ErrNotAuthorized
	}

	user, err := model.User{}.FindUserById(ctx.Request.Context(), c.db, id)
	if err != nil {
		return ErrCredentialInvalid
	}

	response := checkDataCustomer(user, userRequest)

	err = model.User{}.UpdateCustomer(ctx.Request.Context(), c.db, response)
	if err != nil {
		return ErrCredentialInvalid
	}

	return nil
}

func (c *controller) DeleteCustomer(ctx *gin.Context, token string, id int) error {

	_, err := c.verifyToken(token)
	if err != nil {
		return ErrNotAuthorized
	}

	err = model.User{}.DeleteCustomer(ctx.Request.Context(), c.db, id)
	if err != nil {
		return ErrCredentialInvalid
	}

	return nil
}

func checkDataCustomer(user model.User, userRequest model.User) model.User {

	if userRequest.Fullname == "" {
		userRequest.Fullname = user.Fullname
	}

	if userRequest.Email == "" {
		userRequest.Email = user.Email
	}

	if userRequest.Phone == "" {
		userRequest.Phone = user.Phone
	}

	userRequest.ID = user.ID

	return userRequest
}

func (c *controller) GetDetailCustomer(ctx *gin.Context, token string) ([]model.ResponseUserDetail, error) {

	var responseUserDetail []model.ResponseUserDetail

	payload, err := c.verifyToken(token)
	id := payload.UserId
	if err != nil {
		return responseUserDetail, ErrNotAuthorized
	}

	responseUserDetail, err = model.User{}.GetDetailCustomer(ctx.Request.Context(), c.db, int(id))
	if err != nil {
		return nil, err
	}

	return responseUserDetail, nil
}
