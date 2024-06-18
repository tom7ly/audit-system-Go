package handler

import (
	"audit-system/internal/model"
	"audit-system/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var accountService *service.AccountService

func InitAccountHandler(as *service.AccountService) {
	accountService = as
}

type AccountUriRequest struct {
	Email string `uri:"email" binding:"required,email"`
}
type CreateAccountRequest struct {
	Balance float64 `json:"balance" binding:"required"`
}

func CreateAccount(c *gin.Context) {
	email := c.Param("email")
	var req CreateAccountRequest
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleRequestParsingError(c, err)
		return
	}

	account := model.Account{
		Balance:          req.Balance,
		LastTransferTime: time.Now(),
	}
	createdAccount, err := accountService.CreateAccount(c.Request.Context(), account, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, createdAccount)
}

func GetAccountsByEmail(c *gin.Context) {
	var uri AccountUriRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		handleRequestParsingError(c, err)
		return
	}

	accounts, err := accountService.GetAccountsByEmail(c.Request.Context(), uri.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

type AccountByIDUriRequest struct {
	Email string `uri:"email" binding:"required,email"`
	ID    int    `uri:"account_id" binding:"required"`
}

func GetAccountById(c *gin.Context) {
	email := c.Param("email")
	accountID, err := strconv.Atoi(c.Param("accountID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := accountService.GetAccountWithTransactions(c.Request.Context(), email, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

func UpdateAccount(c *gin.Context) {
	email := c.Param("email")
	accountID, err := strconv.Atoi(c.Param("accountID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var account model.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		handleRequestParsingError(c, err)
		return
	}
	updatedAccount, err := accountService.GetAccountByID(c.Request.Context(), email, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedAccount)
}
func DeleteAccount(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("accountID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}
	if err := accountService.DeleteAccount(c.Request.Context(), accountID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
