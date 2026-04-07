package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rdee29/money-keeper/config"
	"github.com/rdee29/money-keeper/internal/handler"
	"github.com/rdee29/money-keeper/internal/model"
	"github.com/rdee29/money-keeper/internal/middleware"
)

func main() {
	config.ConnectDB()

	config.DB.AutoMigrate(&model.User{})
	config.DB.AutoMigrate(&model.Transaction{})

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/profile", middleware.AuthMiddleware(), func(c *gin.Context) {
		userID, _ := c.Get("user_id")

		c.JSON(200, gin.H{
			"message" : "this is your profile",
			"user_id" : userID,
		})
	})

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.POST("/transactions", middleware.AuthMiddleware(), handler.CreateTransaction)
	r.GET("/transactions", middleware.AuthMiddleware(), handler.GetTransactions)


	
	r.Run(":8080")
}
