package main

import (
	"log"
	"net/http"
	"rbac/pkg/db"
	"rbac/pkg/kubernetes"
	"rbac/pkg/server"
)

func main() {
	// Initialize the database
	db.InitDB("db.db")

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewClientset()
	if err != nil {
		log.Fatalf("Error creating Kubernetes clientset: %v", err)
	}

	// Load server configuration
	serverConfig := server.NewConfig()

	// Create and start the server
	srv := server.NewServer(clientset, serverConfig)
	log.Printf("Starting server on port %s", serverConfig.Port)

	if serverConfig.IsDevMode {
		// In development mode, use HTTP
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	} else {
		// In production mode, use HTTPS with certificates managed by the handler
		if err := server.StartServer(srv, serverConfig); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}
}