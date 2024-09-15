package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rbac/pkg/db"
	"rbac/pkg/kubernetes"
	"rbac/pkg/server"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	// Initialize the logger
	utils.InitLogger()
	defer utils.Logger.Sync()

	// Initialize the database
	db.InitDB("db.db")

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewClientset()
	if err != nil {
		utils.Logger.Fatal("Error creating Kubernetes clientset", zap.Error(err))
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
		utils.Logger.Fatal("Error retrieving certificates from database", zap.Error(err))
	}

	// Start server
	go func() {
		if len(certData) > 0 && len(keyData) > 0 {
			// Start server with SSL
			certFile := "/tmp/tls.crt"
			keyFile := "/tmp/tls.key"
			if err := os.WriteFile(certFile, certData, 0644); err != nil {
				utils.Logger.Fatal("Error writing cert file", zap.Error(err))
			}
			if err := os.WriteFile(keyFile, keyData, 0644); err != nil {
				utils.Logger.Fatal("Error writing key file", zap.Error(err))
			}
			if err := e.StartTLS(":"+serverConfig.Port, certFile, keyFile); err != nil && err != http.ErrServerClosed {
				utils.Logger.Fatal("Shutting down the server", zap.Error(err))
			}
			return
		}
		// Start server without SSL
		if err := e.Start(":" + serverConfig.Port); err != nil && err != http.ErrServerClosed {
			utils.Logger.Fatal("Shutting down the server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		utils.Logger.Fatal("Error during server shutdown", zap.Error(err))
	}
}
