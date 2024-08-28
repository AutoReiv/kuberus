package main

import (
	"log"
	"net/http"
	"rbac/pkg/kubernetes"
	"rbac/pkg/server"
)

func main() {
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
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}