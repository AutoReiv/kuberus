package handlers

import (
	"os"
	"path/filepath"
)

// GetKubeConfig returns the path to the kubeconfig file
func GetKubeConfig() string {
	return filepath.Join(homeDir(), ".kube", "config")
}

// homeDir returns the home directory of the user
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
