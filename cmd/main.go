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

func main() {
	var container = service.GetContainer()
	go scheduleAuditLogCleanup(container, 1*time.Minute, 30*time.Second)

	r := gin.Default()
	router.SetupRoutes(r)
	r.Run()
}
func scheduleAuditLogCleanup(container *service.Container, ttl, interval time.Duration) {

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		ctx := context.WithValue(context.Background(), utils.AuditContextKey, true)
		deletes, err := container.AuditLogService.DeleteOldAuditLogs(ctx, ttl)
		if err != nil {
			log.Printf("Failed to delete old audit logs: %v", err)
		} else if deletes > 0 {
			log.Println("Deleted", deletes, "old audit logs")
		} else {
			log.Println("No old audit logs to delete")
		}
	}
}
