package handler

import (
	// "net/http"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rdee29/money-keeper/config"
	"github.com/rdee29/money-keeper/internal/model"
)

type CreateTransactionRequest struct {
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
}

func CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// get user from token
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	req.Type = strings.TrimSpace(strings.ToLower(req.Type))

	if req.Type != model.TypeExpense && req.Type != model.TypeIncome {
		c.JSON(400, gin.H{
			"error" : "type must be either 'income' or 'expense'",
		})
		return
	}
	
	transaction := model.Transaction{
		ID:          uuid.New(),
		UserID:      userID,
		Amount:      req.Amount,
		Type:        req.Type,
		Description: req.Description,
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to create transaction",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "transaction created",
	})
}

func GetTransactions(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var transactions []model.Transaction

	// get query param
	typeQuery := c.Query("type")

	query := config.DB.Where("user_id = ?", userID)

	// filter type
	if typeQuery != "" {
		typeQuery = strings.ToLower(typeQuery)

		if typeQuery != model.TypeIncome && typeQuery != model.TypeExpense {
			c.JSON(400, gin.H{
				"error": "invalid type filter",
			})
			return
		}

		query = query.Where("type = ?", typeQuery)
	}

	if err := query.Find(&transactions).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch transactions",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": transactions,
	})
}

func GetSummary(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var totalIncome float64
	var totalExpense float64

	// summary of income
	config.DB.Model(&model.Transaction{}).
		Where("user_id = ? AND type = ?", userID, model.TypeIncome).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)

	// summary of expense
	config.DB.Model(&model.Transaction{}).
		Where("user_id = ? AND type = ?", userID, model.TypeExpense).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense)

	balance := totalIncome - totalExpense

	c.JSON(200, gin.H{
		"total_income":  totalIncome,
		"total_expense": totalExpense,
		"balance":       balance,
	})
}