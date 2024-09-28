package handlers

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/db"
	"rbac/pkg/utils"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// UploadCertsHandler handles the uploading, updating, and deleting of TLS certificates.
type UploadCertsHandler struct {
	Clientset *kubernetes.Clientset
}

func NewUploadCertsHandler(clientset *kubernetes.Clientset) *UploadCertsHandler {
	return &UploadCertsHandler{Clientset: clientset}
}

func (h *UploadCertsHandler) ServeHTTP(c echo.Context) error {
	username, _ := c.Get("username").(string)
	isAdmin, ok := c.Get("isAdmin").(bool)
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Unable to determine admin status")
	}

	if !isAdmin {
		utils.Logger.Warn("Forbidden: user is not an admin", zap.String("username", username))
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	switch c.Request().Method {
	case http.MethodPost:
		return h.handleUpload(c)
	case http.MethodPut:
		return h.handleUpdate(c)
	case http.MethodDelete:
		return h.handleDelete(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed, "Method not allowed")
	}
}
func (h *UploadCertsHandler) handleUpload(c echo.Context) error {
	// Ensure the user is an admin
	username, ok := c.Get("username").(string)
	if !ok || !auth.IsAdmin(username) {
		utils.Logger.Warn("Forbidden: user is not an admin", zap.String("username", username))
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	// Parse the multipart form
	if err := c.Request().ParseMultipartForm(10 << 20); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to parse form", err, "Failed to parse multipart form")
	}

	// Retrieve the certificate and key files
	certFile, _, err := c.Request().FormFile("certFile")
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to retrieve cert file", err, "Failed to retrieve cert file from form")
	}
	defer certFile.Close()

	keyFile, _, err := c.Request().FormFile("keyFile")
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to retrieve key file", err, "Failed to retrieve key file from form")
	}
	defer keyFile.Close()

	// Read the certificate and key files
	certData, err := io.ReadAll(certFile)
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to read cert file", err, "Failed to read cert file")
	}

	keyData, err := io.ReadAll(keyFile)
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to read key file", err, "Failed to read key file")
	}

	// Validate the certificate and key files
	if err := isValidCertAndKey(certData, keyData); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid certificate or key", err, "Invalid certificate or key")
	}

	// Store the certificates in the database
	_, err = db.DB.Exec("INSERT INTO certificates (cert, key) VALUES (?, ?)", certData, keyData)
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to store certificates", err, "Failed to store certificates in database")
	}

	utils.Logger.Info("Certificates uploaded successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "upload_certs", "N/A", "N/A")

	// Respond with success
	utils.WriteJSON(c.Response(), map[string]string{"message": "Certificates uploaded successfully. Deployment will restart to apply changes. This may take a few moments."})

	// Trigger a rolling restart of the deployment
	go h.triggerRollingRestart()
	return nil
}

func (h *UploadCertsHandler) handleUpdate(c echo.Context) error {
	// Ensure the user is an admin
	username, ok := c.Get("username").(string)
	if !ok || !auth.IsAdmin(username) {
		utils.Logger.Warn("Forbidden: user is not an admin", zap.String("username", username))
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	// Parse the multipart form
	if err := c.Request().ParseMultipartForm(10 << 20); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to parse form", err, "Failed to parse multipart form")
	}

	// Retrieve the certificate and key files
	certFile, _, err := c.Request().FormFile("certFile")
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to retrieve cert file", err, "Failed to retrieve cert file from form")
	}
	defer certFile.Close()

	keyFile, _, err := c.Request().FormFile("keyFile")
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to retrieve key file", err, "Failed to retrieve key file from form")
	}
	defer keyFile.Close()

	// Read the certificate and key files
	certData, err := io.ReadAll(certFile)
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to read cert file", err, "Failed to read cert file")
	}

	keyData, err := io.ReadAll(keyFile)
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to read key file", err, "Failed to read key file")
	}

	// Validate the certificate and key files
	if err := isValidCertAndKey(certData, keyData); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid certificate or key", err, "Invalid certificate or key")
	}

	// Update the certificates in the database
	_, err = db.DB.Exec("UPDATE certificates SET cert = ?, key = ? WHERE id = (SELECT id FROM certificates ORDER BY created_at DESC LIMIT 1)", certData, keyData)
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to update certificates", err, "Failed to update certificates in database")
	}

	utils.Logger.Info("Certificates updated successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "update_certs", "N/A", "N/A")

	// Respond with success
	utils.WriteJSON(c.Response(), map[string]string{"message": "Certificates updated successfully. Deployment will restart to apply changes. This may take a few moments."})

	// Trigger a rolling restart of the deployment
	go h.triggerRollingRestart()
	return nil
}

func (h *UploadCertsHandler) handleDelete(c echo.Context) error {
	// Ensure the user is an admin
	username, ok := c.Get("username").(string)
	if !ok || !auth.IsAdmin(username) {
		utils.Logger.Warn("Forbidden: user is not an admin", zap.String("username", username))
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	// Delete the certificates from the database
	_, err := db.DB.Exec("DELETE FROM certificates WHERE id = (SELECT id FROM certificates ORDER BY created_at DESC LIMIT 1)")
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to delete certificates", err, "Failed to delete certificates from database")
	}

	utils.Logger.Info("Certificates deleted successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "delete_certs", "N/A", "N/A")

	// Respond with success
	utils.WriteJSON(c.Response(), map[string]string{"message": "Certificates deleted successfully. Deployment will restart to apply changes. This may take a few moments."})

	// Trigger a rolling restart of the deployment
	go h.triggerRollingRestart()
	return nil
}

func isValidCertAndKey(certData, keyData []byte) error {
	_, err := tls.X509KeyPair(certData, keyData)
	return err
}

func (h *UploadCertsHandler) triggerRollingRestart() {
	deploymentsClient := h.Clientset.AppsV1().Deployments("default")
	deployment, err := deploymentsClient.Get(context.TODO(), "your-deployment-name", metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Failed to get deployment", zap.Error(err))
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
		utils.Logger.Error("Failed to update deployment", zap.Error(err))
	}
}