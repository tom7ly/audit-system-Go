package router

import (
	"audit-system/internal/database"
	"audit-system/internal/handler"
	"audit-system/internal/repository"
	"audit-system/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, auditLogRepo *repository.AuditLogRepository) {
	// Initialize Repositories
	userRepository := repository.NewUserRepository(database.Client)
	accountRepository := repository.NewAccountRepository(database.Client)
	transactionRepository := repository.NewTransactionRepository(database.Client)

	// Initialize Services
	userService := service.NewUserService(userRepository)
	accountService := service.NewAccountService(accountRepository)
	transactionService := service.NewTransactionService(transactionRepository)
	auditLogService := service.NewAuditLogService(auditLogRepo)

	// Initialize Handlers
	handler.InitUserHandler(userService)
	handler.InitAccountHandler(accountService)
	handler.InitTransactionHandler(transactionService)
	handler.InitAuditLogHandler(auditLogService)

	// User Routes
	userGroup := r.Group("/users")
	{
		userGroup.POST("/", handler.CreateUser)
		userGroup.GET("/", handler.GetUsers)
		userGroup.GET("/:email", handler.GetUserByEmail)
		userGroup.PUT("/:email", handler.UpdateUser)
		userGroup.GET("/:email/accounts", handler.GetAccountsByEmail)
		userGroup.POST("/:email/accounts", handler.CreateAccount)
		userGroup.GET("/:email/accounts/:accountID", handler.GetAccountById)
		userGroup.PUT("/:email/accounts/:accountID", handler.UpdateAccount)
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
