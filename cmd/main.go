package main

import (
	"log"
	"net/http"
	"os"
	"rbac/pkg/kubernetes"
	"rbac/pkg/server"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Set Gin mode based on the environment variable
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode // Default to debug mode if not set
	}
	gin.SetMode(ginMode)

	// Debug prints
	log.Printf("DEV_MODE: %s", os.Getenv("DEV_MODE"))
	log.Printf("CERT_FILE: %s", os.Getenv("CERT_FILE"))
	log.Printf("KEY_FILE: %s", os.Getenv("KEY_FILE"))
	log.Printf("PORT: %s", os.Getenv("PORT"))
	log.Printf("GIN: %s", ginMode)

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewClientset()
	if err != nil {
		log.Fatalf("Error creating Kubernetes clientset: %v", err)
	}

	// Load server configuration
	serverConfig := server.NewConfig()

	// Set the session secret for Goth
	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	// Read certificate and key file paths from environment variables
	certFile := os.Getenv("CERT_FILE")
	keyFile := os.Getenv("KEY_FILE")

	// Create and start the server
	srv := server.NewServer(clientset, serverConfig)
	log.Printf("Starting server on port %s", serverConfig.Port)

	if serverConfig.IsDevMode || certFile == "" || keyFile == "" {
		// In development mode or if certs are not provided, use HTTP
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	} else {
		// In production mode, use HTTPS
		if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}
}
