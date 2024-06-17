package router

import (
	"audit-system/internal/database"
	"audit-system/internal/handler"
	"audit-system/internal/repository"
	"audit-system/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	userRepository := repository.NewUserRepository(database.Client)
	userService := service.NewUserService(userRepository)
	handler.InitUserHandler(userService)

	accountRepository := repository.NewAccountRepository(database.Client)
	accountService := service.NewAccountService(accountRepository)
	handler.InitAccountHandler(accountService)

	transactionRepository := repository.NewTransactionRepository(database.Client)
	transactionService := service.NewTransactionService(transactionRepository)
	handler.InitTransactionHandler(transactionService)

	auditLogRepository := repository.NewAuditLogRepository(database.Client)
	auditLogService := service.NewAuditLogService(auditLogRepository)
	handler.InitAuditLogHandler(auditLogService)

	userGroup := r.Group("/users")
	{
		userGroup.POST("/", handler.CreateUser)
		userGroup.GET("/", handler.GetUsers)
		userGroup.GET("/:email", handler.GetUserByEmail)
		userGroup.PUT("/:email", handler.UpdateUser)
		userGroup.GET("/:email/accounts", handler.GetAccountsByEmail)
		userGroup.GET("/:email/accounts/:accountID", handler.GetAccountById)
		userGroup.POST("/:email/accounts", handler.CreateAccount)
	}

	accountGroup := r.Group("/accounts")
	{
		accountGroup.GET("/:email", handler.GetAccountsByEmail)
		accountGroup.GET("/:email/:accountID", handler.GetAccountById)
		accountGroup.GET("/:email/:accountID/transactions/inbound", func(c *gin.Context) {
			c.Set("transaction_type", "inbound")
			handler.GetTransactions(c)
		})
		accountGroup.GET("/:email/:accountID/transactions/outbound", func(c *gin.Context) {
			c.Set("transaction_type", "outbound")
			handler.GetTransactions(c)
		})
		accountGroup.GET("/:email/:accountID/transactions", func(c *gin.Context) {
			c.Set("transaction_type", "")
			handler.GetTransactions(c)
		})
	}

	transactionGroup := r.Group("/transactions")
	{
		transactionGroup.POST("/:email/:accountID", handler.CreateTransaction)
		transactionGroup.GET("/:email/:accountID/inbound", func(c *gin.Context) {
			c.Set("transaction_type", "inbound")
			handler.GetTransactions(c)
		})
		transactionGroup.GET("/:email/:accountID/outbound", func(c *gin.Context) {
			c.Set("transaction_type", "outbound")
			handler.GetTransactions(c)
		})
		transactionGroup.GET("/:email/:accountID", func(c *gin.Context) {
			c.Set("transaction_type", "")
			handler.GetTransactions(c)
		})
	}

	auditLogGroup := r.Group("/auditlogs")
	{
		auditLogGroup.GET("/", handler.GetAllAuditLogs)
		auditLogGroup.GET("/:email", handler.GetAuditLogsByEmail)
	}

}
