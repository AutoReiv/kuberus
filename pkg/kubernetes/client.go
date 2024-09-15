package kubernetes

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"rbac/pkg/utils"
	"go.uber.org/zap"
)

// NewClientset initializes and returns a Kubernetes clientset
func NewClientset() (*kubernetes.Clientset, error) {
	kubeconfig := filepath.Join(homeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		utils.Logger.Error("Error building Kubernetes config", zap.Error(err))
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		utils.Logger.Error("Error creating Kubernetes clientset", zap.Error(err))
		return nil, err
	}

	utils.Logger.Info("Kubernetes clientset created successfully")
	return clientset, nil
}

// homeDir returns the home directory of the user
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
