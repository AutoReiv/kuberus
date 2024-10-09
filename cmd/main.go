package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rbac/pkg/kubernetes"
	"rbac/pkg/server"

	"github.com/labstack/echo/v4"
	"github.com/rs/cors"
)

func main() {
	// Create Kubernetes clientset
	clientset, err := kubernetes.NewClientset()
	if err != nil {
		panic("Error creating Kubernetes clientset: " + err.Error())
	}

	// Create Echo instance
	e := echo.New()

	// CORS
	e.Use(echo.WrapMiddleware(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler))

	// Load server configuration
	serverConfig := server.NewConfig()

	// Register routes
	server.RegisterRoutes(e, clientset, serverConfig)

	// Start server
	go func() {
		if err := e.Start(":" + serverConfig.Port); err != nil && err != http.ErrServerClosed {
			panic("Shutting down the server: " + err.Error())
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		panic("Error during server shutdown: " + err.Error())
	}
}
