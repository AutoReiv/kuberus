package main

import (
	"log"
	"net/http"
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
		ClientID:     "b7Ql1vVr7bP1Tw1b7DQj5eROQwIY6cGG",
		ClientSecret: "Pdn94ZBVt7RU3eNsFZ25XaxES_w7_4VThC6Dom6U4O7n26YHfgR4OzujFSG22Xl6",
		IssuerURL:    "https://dev-ooomxzist3l3qxwf.us.auth0.com/",
		RedirectURL:  "http://localhost:8080/callback",
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
