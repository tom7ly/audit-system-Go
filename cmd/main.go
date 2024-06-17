package main

import (
	"audit-system/internal/database"
	"audit-system/internal/repository"
	"audit-system/internal/router"
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	database.Init()
	defer database.Close()

	// Initialize repositories
	auditLogRepo := repository.NewAuditLogRepository(database.Client)

	// Schedule the cleanup job
	go scheduleAuditLogCleanup(auditLogRepo, 10*time.Minute, 1*time.Minute)

	// Initialize the router and set up routes
	r := gin.Default()
	router.SetupRoutes(r, auditLogRepo)
	r.Run() // Start the server
}

func scheduleAuditLogCleanup(repo *repository.AuditLogRepository, ttl, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		err := repo.DeleteOldAuditLogs(context.Background(), ttl)
		if err != nil {
			log.Printf("Failed to delete old audit logs: %v", err)
		} else {
			log.Println("Old audit logs deleted successfully")
		}
	}
}
