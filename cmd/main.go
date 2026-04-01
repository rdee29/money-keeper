package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rdee29/money-keeper/config"
	"github.com/rdee29/money-keeper/internal/model"
)

func main() {
	config.ConnectDB()

	config.DB.AutoMigrate(&model.User{})

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}
