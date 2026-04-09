package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods: 	  []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
	r.GET("/summary", middleware.AuthMiddleware(), handler.GetSummary)

	r.Run(":8080")
}
