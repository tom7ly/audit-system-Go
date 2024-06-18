package service

import (
	"audit-system/internal/repository"
	"audit-system/internal/utils"
	"sync"
)

type Container struct {
	DBService             *DBService
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
		instance.initialize("")
	}
	return instance
}

func GetTestContainer(dsn string) *Container {
	if instance == nil {
		instance = &Container{}
		instance.initialize(dsn)
	}
	return instance
}

func (c *Container) initialize(dsn string) {
	c.once.Do(func() {
		c.DBService = GetDBService()
		c.DBService.Init(dsn)

		c.Queue = utils.NewQueue(100, 10)
		c.UserRepository = repository.NewUserRepository(c.DBService.Client(), c.Queue)
		c.AccountRepository = repository.NewAccountRepository(c.DBService.Client(), c.Queue)
		c.TransactionRepository = repository.NewTransactionRepository(c.DBService.Client(), c.Queue)
		c.AuditLogRepository = repository.NewAuditLogRepository(c.DBService.Client(), c.Queue)
		c.UserService = newUserService(c.UserRepository, c.AccountRepository)
		c.AccountService = newAccountService(c.AccountRepository)
		c.TransactionService = NewTransactionService(c.TransactionRepository)
		c.AuditLogService = NewAuditLogService(c.AuditLogRepository)
	})
}

func (c *Container) Shutdown() {
	c.DBService.Close()
}
