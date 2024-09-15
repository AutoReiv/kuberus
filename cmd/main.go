package main

import (
	"context"
	"database/sql"
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

	// Load server configuration
	serverConfig := server.NewConfig()

	// Register routes
	server.RegisterRoutes(e, clientset, serverConfig)

	// Check if SSL certificates are available in the database
	var certData, keyData []byte
	err = db.DB.QueryRow("SELECT cert, key FROM certificates ORDER BY created_at DESC LIMIT 1").Scan(&certData, &keyData)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("Error retrieving certificates from database: %v", err)
	}

	// Start server
	go func() {
		if len(certData) > 0 && len(keyData) > 0 {
			// Start server with SSL
			certFile := "/tmp/tls.crt"
			keyFile := "/tmp/tls.key"
			if err := os.WriteFile(certFile, certData, 0644); err != nil {
				log.Fatalf("Error writing cert file: %v", err)
			}
			if err := os.WriteFile(keyFile, keyData, 0644); err != nil {
				log.Fatalf("Error writing key file: %v", err)
			}
			if err := e.StartTLS(":"+serverConfig.Port, certFile, keyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Shutting down the server: %v", err)
			}
			return
		}
		// Start server without SSL
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
