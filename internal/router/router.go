package router

import (
	"audit-system/internal/handler"
	"audit-system/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Initialize Repositories
	container := service.GetContainer()
	// Initialize Handlers
	handler.InitUserHandler(container.UserService)
	handler.InitAccountHandler(container.AccountService)
	handler.InitTransactionHandler(container.TransactionService)
	handler.InitAuditLogHandler(container.AuditLogService)

	// User Routes
	userGroup := r.Group("/users")
	{
		userGroup.POST("/", handler.CreateUser)
		userGroup.GET("/", handler.GetUsers)
		userGroup.GET("/:email", handler.GetUserByEmail)
		userGroup.PUT("/:email", handler.UpdateUser)
		userGroup.DELETE("/:email", handler.DeleteUser) // New endpoint to delete user
		userGroup.GET("/:email/accounts", handler.GetAccountsByEmail)
		userGroup.POST("/:email/accounts", handler.CreateAccount)
		userGroup.GET("/:email/accounts/:accountID", handler.GetAccountById)
		userGroup.PUT("/:email/accounts/:accountID", handler.UpdateAccount)
		userGroup.DELETE("/:email/accounts/:accountID", handler.DeleteAccount) // New endpoint to delete account by user and account ID
	}

	// Account Routes
	accountGroup := r.Group("/accounts")
	{
		accountGroup.POST("/:email", handler.CreateAccount)
		accountGroup.GET("/:email", handler.GetAccountsByEmail)
		accountGroup.GET("/:email/:accountID", handler.GetAccountById)
		accountGroup.GET("/:email/:accountID/transactions", setTransactionType("", handler.GetTransactions))
		accountGroup.GET("/:email/:accountID/transactions/inbound", setTransactionType("inbound", handler.GetTransactions))
		accountGroup.GET("/:email/:accountID/transactions/outbound", setTransactionType("outbound", handler.GetTransactions))
		accountGroup.DELETE("/:accountID", handler.DeleteAccount) // New endpoint to delete account by account ID only
	}

	// Transaction Routes
	transactionGroup := r.Group("/transactions")
	{
		transactionGroup.POST("/:email/:accountID", handler.CreateTransaction)
		transactionGroup.GET("/:email/:accountID/inbound", setTransactionType("inbound", handler.GetTransactions))
		transactionGroup.GET("/:email/:accountID/outbound", setTransactionType("outbound", handler.GetTransactions))
		transactionGroup.GET("/:email/:accountID", setTransactionType("", handler.GetTransactions))
	}

	// Audit Log Routes
	auditLogGroup := r.Group("/auditlogs")
	{
		auditLogGroup.GET("/", handler.GetAllAuditLogs)
		auditLogGroup.GET("/:email", handler.GetAuditLogsByEmail)
	}
}

// Middleware to set transaction type in context
func setTransactionType(transactionType string, handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("transaction_type", transactionType)
		handlerFunc(c)
	}
}
