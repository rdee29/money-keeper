package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rdee29/money-keeper/config"
	"github.com/rdee29/money-keeper/internal/model"
)

type RegisterRequest struct {
	Name 		string `json:"name"`
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

func Register (c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := model.User {
		ID: uuid.New(),
		Name: req.Name,
		Email: req.Email,
		Password: req.Password,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : "failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message" : "user successfully created",
	})
}