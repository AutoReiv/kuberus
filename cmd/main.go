package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rbac/pkg/db"
	"rbac/pkg/kubernetes"
	"rbac/pkg/server"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize the database
	db.InitDB("db.db")

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewClientset()
	if err != nil {
		log.Fatalf("Error creating Kubernetes clientset: %v", err)
	}

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	// Load server configuration
	serverConfig := server.NewConfig()

	// Register routes
	server.RegisterRoutes(e, clientset, serverConfig)

	// Start server
	go func() {
		if err := e.Start(":" + serverConfig.Port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Shutting down the server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
