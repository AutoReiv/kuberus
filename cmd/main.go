package main

import (
	"log"
	"net/http"
	"os"
	"rbac/pkg/kubernetes"
	"rbac/pkg/server"
	"rbac/pkg/utils"
)

func main() {
	// Create Kubernetes clientset
	clientset, err := kubernetes.NewClientset()
	if err != nil {
		log.Fatalf("Error creating Kubernetes clientset: %v", err)
	}

	// Load server configuration
	serverConfig := server.NewConfig()

	// Paths to the certificate and key files
	certFile := "cert.pem"
	keyFile := "key.pem"

	// Generate self-signed certificates if they do not exist
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Println("Generating self-signed certificates...")
		if err := utils.GenerateSelfSignedCert(certFile, keyFile); err != nil {
			log.Fatalf("Error generating self-signed certificates: %v", err)
		}
	}

	// Create and start the server
	srv := server.NewServer(clientset, serverConfig)
	log.Printf("Starting server on port %s", serverConfig.Port)
	if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}