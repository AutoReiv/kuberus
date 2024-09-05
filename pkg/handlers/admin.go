package handlers

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreateAdminRequest represents the request payload for creating an admin account.
type CreateAdminRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required,eqfield=Password"`
}

// CreateAdminHandler handles the creation of an admin account.
func CreateAdminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Sanitize user input
	username := utils.SanitizeInput(req.Username)
	password := utils.SanitizeInput(req.Password)

	// Validate password strength
	if !utils.IsStrongPassword(password) {
		http.Error(w, "Password does not meet strength requirements", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "Error hashing password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Acquire lock to synchronize access to shared data
	auth.Mu.Lock()
	defer auth.Mu.Unlock()

	// Check if an admin account already exists
	if auth.AdminExists {
		http.Error(w, "Admin account already exists", http.StatusConflict)
		return
	}

	// Store the admin account information
	auth.Users[username] = hashedPassword
	auth.AdminExists = true

	utils.WriteJSON(w, map[string]string{"message": "Admin account created successfully"})
}

// UploadCertsHandler handles the uploading of TLS certificates.
type UploadCertsHandler struct {
	Clientset *kubernetes.Clientset
}

func NewUploadCertsHandler(clientset *kubernetes.Clientset) *UploadCertsHandler {
	return &UploadCertsHandler{Clientset: clientset}
}

func (h *UploadCertsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ensure the user is an admin
	username, ok := r.Context().Value("username").(string)
	if !ok || !auth.IsAdmin(username) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the certificate and key files
	certFile, _, err := r.FormFile("certFile")
	if err != nil {
		http.Error(w, "Failed to retrieve cert file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer certFile.Close()

	keyFile, _, err := r.FormFile("keyFile")
	if err != nil {
		http.Error(w, "Failed to retrieve key file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer keyFile.Close()

	// Read the certificate and key files
	certData, err := io.ReadAll(certFile)
	if err != nil {
		http.Error(w, "Failed to read cert file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	keyData, err := io.ReadAll(keyFile)
	if err != nil {
		http.Error(w, "Failed to read key file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate the certificate and key files
	if err := isValidCertAndKey(certData, keyData); err != nil {
		http.Error(w, "Invalid certificate or key: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create or update the Kubernetes Secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tls-certs",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"tls.crt": certData,
			"tls.key": keyData,
		},
	}

	_, err = h.Clientset.CoreV1().Secrets("default").Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		_, err = h.Clientset.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		if err != nil {
			http.Error(w, "Failed to create/update secret: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Respond with success
	utils.WriteJSON(w, map[string]string{"message": "Certificates uploaded successfully. Deployment will restart to apply changes. This may take a few moments."})

	// Trigger a rolling restart of the deployment
	go h.triggerRollingRestart()
}

func isValidCertAndKey(certData, keyData []byte) error {
	_, err := tls.X509KeyPair(certData, keyData)
	return err
}

func (h *UploadCertsHandler) triggerRollingRestart() {
	deploymentsClient := h.Clientset.AppsV1().Deployments("default")
	deployment, err := deploymentsClient.Get(context.TODO(), "your-deployment-name", metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get deployment: %v", err)
		return
	}

	// Update the deployment to trigger a rolling restart
	annotations := deployment.Spec.Template.Annotations
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
	deployment.Spec.Template.Annotations = annotations

	_, err = deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Failed to update deployment: %v", err)
	}
}
