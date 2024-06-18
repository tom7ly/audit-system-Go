package service

import (
	"audit-system/internal/repository"
	"audit-system/internal/utils"
	"sync"
)

type Container struct {
	Queue                 *utils.Queue
	UserRepository        *repository.UserRepository
	AccountRepository     *repository.AccountRepository
	TransactionRepository *repository.TransactionRepository
	UserService           *UserService
	AccountService        *AccountService
	TransactionService    *TransactionService
	AuditLogService       *AuditLogService
	AuditLogRepository    *repository.AuditLogRepository
	once                  sync.Once
}

var instance *Container

func GetContainer() *Container {
	if instance == nil {
		instance = &Container{}
		instance.initialize()
	}
	return instance
}

func (c *Container) initialize() {
	c.once.Do(func() {
		dbService := GetDBService()
		dbService.Init()

		c.Queue = utils.NewQueue(100, 10)
		c.UserRepository = repository.NewUserRepository(dbService.Client(), c.Queue)
		c.AccountRepository = repository.NewAccountRepository(dbService.Client(), c.Queue)
		c.TransactionRepository = repository.NewTransactionRepository(dbService.Client(), c.Queue)
		c.AuditLogRepository = repository.NewAuditLogRepository(dbService.client, c.Queue)
		c.UserService = NewUserService(c.UserRepository, c.AccountRepository)
		c.AccountService = newAccountService(c.AccountRepository)
		c.TransactionService = NewTransactionService(c.TransactionRepository)
		c.AuditLogService = NewAuditLogService(c.AuditLogRepository)
	})
}

func (c *Container) Shutdown() {
	dbService := GetDBService()
	dbService.Close()
}
