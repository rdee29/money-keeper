package handler

import (
	// "errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rdee29/money-keeper/config"
	"github.com/rdee29/money-keeper/internal/middleware"
	"github.com/rdee29/money-keeper/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to hashed password",
		})
		return
	}

	user := model.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user successfully created",
	})
}

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user model.User

	//find by email
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	//compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(401, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	token, err := middleware.GeneratToken(user.ID.String())
	if err != nil {
		c.JSON(500, gin.H{
			"error" : "failed to generate token",
		})
		return
	}

	c.JSON(200, gin.H{
		"message" : "login successful",
		"token" : token,
	})
}
