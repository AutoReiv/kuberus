package main

import (
	"log"
	"net/http"
	"os"
	"rbac/pkg/auth"
	"rbac/pkg/kubernetes"
	"rbac/pkg/server"
)

func main() {
	// Initialize the Kubernetes client
	clientset, err := kubernetes.NewClientset()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	// Initialize OIDC
	oidcConfig := auth.OIDCConfig{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		IssuerURL:    os.Getenv("ISSUER_URL"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
	}
	if err := auth.InitOIDC(oidcConfig); err != nil {
		log.Fatalf("Failed to initialize OIDC: %v", err)
	}

	// Initialize the server
	srv := server.NewServer(clientset)

	// Apply OIDC Middleware
	http.Handle("/", auth.OIDCMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, authenticated user!"))
	})))

	// Start the server
	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
