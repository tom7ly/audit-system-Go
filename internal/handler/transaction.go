package handler

import (
	"audit-system/internal/service"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var transactionService *service.TransactionService

func InitTransactionHandler(ts *service.TransactionService) {
	transactionService = ts
}

type CreateTransactionRequest struct {
	ToAccountID int     `json:"to_account_id"`
	Amount      float64 `json:"amount"`
}

func CreateTransaction(c *gin.Context) {
	email := c.Param("email")
	fromAccountID, err := strconv.Atoi(c.Param("accountID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from account ID"})
		return
	}

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := transactionService.CreateTransaction(context.Background(), email, fromAccountID, req.ToAccountID, req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "transaction created successfully"})
}

func GetTransactions(c *gin.Context) {
	email := c.Param("email")
	accountID, err := strconv.Atoi(c.Param("accountID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}
	transactionType := c.Param("transaction_type")
	if transactionType == "" {
		transactionType = c.GetString("transaction_type")
	}
	transactions, err := transactionService.GetTransactions(context.Background(), email, accountID, transactionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}
