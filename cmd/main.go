package main

import (
	"audit-system/internal/router"
	"audit-system/internal/service"
	"audit-system/internal/utils"
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

var container = *service.GetContainer()

func main() {
	// Initialize database connection

	// Schedule the cleanup job
	go scheduleAuditLogCleanup(1*time.Minute, 30*time.Second)

	// Initialize the router and set up routes
	r := gin.Default()
	router.SetupRoutes(r)
	r.Run() // Start the server
}

func scheduleAuditLogCleanup(ttl, interval time.Duration) {
	repo := container.AuditLogRepository
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		ctx := context.WithValue(context.Background(), utils.AuditContextKey, true)
		deletes, err := repo.DeleteOldAuditLogs(ctx, ttl)
		if err != nil {
			log.Printf("Failed to delete old audit logs: %v", err)
		} else if deletes > 0 {
			log.Println("Deleted", deletes, "old audit logs")
		} else {
			log.Println("No old audit logs to delete")
		}
	}
}
