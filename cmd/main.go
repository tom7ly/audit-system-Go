package main

import (
	"audit-system/internal/database"
	"audit-system/internal/middleware"
	"audit-system/internal/router"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	defer database.Close()

	r := gin.Default()
	r.Use(middleware.UserMiddleware())

	router.SetupRoutes(r)
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("failed to run server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}
