package main

import (
	"log"
	"net/http"
	"os"
	"rbac/pkg/db"
	"rbac/pkg/kubernetes"
	"rbac/pkg/server"
)

func main() {
	// Load environment variables from .env file
	// Commented out the code that loads the .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	// Initialize the database
	db.InitDB("db.db")

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewClientset()
	if err != nil {
		log.Fatalf("Error creating Kubernetes clientset: %v", err)
	}

	// Load server configuration
	serverConfig := server.NewConfig()

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
