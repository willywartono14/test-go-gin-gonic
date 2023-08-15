package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/willywartono14/test-go-gin-gonic/model"
)

type ItemController interface {
	GetAllItem(ctx *gin.Context, token string) ([]model.ResponseItem, error)
}

func (c *controller) GetAllItem(ctx *gin.Context, token string) ([]model.ResponseItem, error) {

	var responseItem []model.ResponseItem

	_, err := c.verifyToken(token)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	item, err := model.Item{}.GetAllItem(ctx.Request.Context(), c.db)
	if err != nil {
		return nil, ErrCredentialInvalid
	}

	if len(item) > 0 {
		for x := range item {
			responseItem = append(responseItem, model.ResponseItem{
				ID:        item[x].ID,
				ItemName:  item[x].ItemName,
				ItemPrice: item[x].ItemPrice,
				ItemStock: item[x].ItemStock,
			})
		}
	}

	return responseItem, nil
}
