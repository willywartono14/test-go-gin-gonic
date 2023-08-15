package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/willywartono14/test-go-gin-gonic/config"
	"github.com/willywartono14/test-go-gin-gonic/controller"
	"github.com/willywartono14/test-go-gin-gonic/database"
	"github.com/willywartono14/test-go-gin-gonic/view"
)

func main() {
	router := gin.Default()

	err := config.Init("config.yaml")
	if err != nil {
		log.Fatalf("[CONFIG] Failed to initialize config: %v", err)
	}

	db := database.Init()

	c := controller.NewController(db)

	view.NewAuthView(c, router)

	router.Run(":8080")
}
