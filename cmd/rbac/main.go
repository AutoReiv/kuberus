package main

import (
	"log"
	"rbac/pkg/kubernetes"
	"rbac/pkg/server"
)

func main() {
	// Initialize the Kubernetes client
	clientset, err := kubernetes.NewClientset()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	// Initialize the server
	srv := server.NewServer(clientset)

	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
