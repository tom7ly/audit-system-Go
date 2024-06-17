package database

import (
	"audit-system/ent"
	"audit-system/internal/hooks"
	"audit-system/internal/repository"
	"context"
	"log"

	"audit-system/internal/service"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

var Client *ent.Client

func Init() {
	var err error
	Client, err = ent.Open("sqlite3", "file:mydatabase.db?_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	auditLogRepo := repository.NewAuditLogRepository(Client)
	auditLogService := service.NewAuditLogService(auditLogRepo)

	// Register the audit log hook
	Client.Use(hooks.AuditLogHook(auditLogService))

	if err := Client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}

func Close() {
	if Client != nil {
		Client.Close()
	}
}
