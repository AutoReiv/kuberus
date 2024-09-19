package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var baseURL string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "mycli",
		Short: "A CLI tool for interacting with the RBAC app",
		Long:  `A CLI tool for logging in and configuring kubectl to interact with the RBAC app.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Please use a subcommand. For example, 'mycli login --url http://localhost:3000'")
		},
	}

	rootCmd.PersistentFlags().StringVarP(&baseURL, "url", "u", "", "Base URL of the RBAC app")
	rootCmd.MarkPersistentFlagRequired("url")
	rootCmd.AddCommand(loginCmd)
	rootCmd.Execute()
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the application and configure kubectl",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting login process...")

		// Open the default web browser and redirect to the login page
		loginURL := fmt.Sprintf("%s/auth/login", baseURL)
		fmt.Println("Opening browser to URL:", loginURL)
		err := openBrowser(loginURL)
		if err != nil {
			fmt.Println("Error opening browser:", err)
			os.Exit(1)
		}

		// Wait for the user to complete the login process and retrieve the token
		token, err := waitForToken()
		if err != nil {
			fmt.Println("Error retrieving token:", err)
			os.Exit(1)
		}

		// Save the token to a file
		tokenFile := filepath.Join(os.Getenv("HOME"), ".kube", "token")
		err = os.WriteFile(tokenFile, []byte(token), 0600)
		if err != nil {
			fmt.Println("Error writing token to file:", err)
			os.Exit(1)
		}
		fmt.Println("Token saved to", tokenFile)

		// Fetch accessible clusters
		clusters, err := fetchAccessibleClusters(token)
		if err != nil {
			fmt.Println("Error fetching accessible clusters:", err)
			os.Exit(1)
		}

		// Configure kubectl for each accessible cluster
		for _, cluster := range clusters {
			err := configureKubectl(cluster)
			if err != nil {
				fmt.Println("Error configuring kubectl for cluster", cluster, ":", err)
				os.Exit(1)
			}
		}
	},
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
	default:
		return fmt.Errorf("unsupported platform")
	}

	return exec.Command(cmd, args...).Start()
}

func waitForToken() (string, error) {
	fmt.Println("Waiting for token...")
	// Poll the server for the token
	for {
		resp, err := http.Get(fmt.Sprintf("%s/auth/token", baseURL))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var result map[string]string
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				return "", err
			}
			return result["token"], nil
		}

		time.Sleep(2 * time.Second)
	}
}

func fetchAccessibleClusters(token string) ([]string, error) {
	fmt.Println("Fetching accessible clusters...")
	req, err := http.NewRequest("GET", baseURL+"/api/clusters", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch clusters: %s", string(body))
	}

	var clusters []string
	err = json.NewDecoder(resp.Body).Decode(&clusters)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

func configureKubectl(cluster string) error {
	fmt.Println("Configuring kubectl for cluster:", cluster)
	// Save the cluster name to a file
	clusterFile := filepath.Join(os.Getenv("HOME"), ".kube", "cluster")
	err := os.WriteFile(clusterFile, []byte(cluster), 0600)
	if err != nil {
		return fmt.Errorf("error writing cluster to file: %w", err)
	}

	// Create the authentication plugin script
	pluginScript := filepath.Join(os.Getenv("HOME"), ".kube", "kubectl-auth-plugin.sh")
	scriptContent := fmt.Sprintf(`#!/bin/bash
TOKEN=$(cat %s)
CLUSTER=$(cat %s)
cat <<EOF
{
  "apiVersion": "client.authentication.k8s.io/v1beta1",
  "kind": "ExecCredential",
  "status": {
    "token": "$TOKEN"
  }
}
EOF
`, filepath.Join(os.Getenv("HOME"), ".kube", "token"), clusterFile)

	err = os.WriteFile(pluginScript, []byte(scriptContent), 0755)
	if err != nil {
		return fmt.Errorf("error writing plugin script to file: %w", err)
	}
	fmt.Println("Authentication plugin script created at", pluginScript)

	// Update kubeconfig to use the authentication plugin
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	kubeconfigContent := fmt.Sprintf(`
apiVersion: v1
clusters:
- cluster:
    server: %s/proxy/%s
  name: %s
contexts:
- context:
    cluster: %s
    user: your-user
  name: %s-context
current-context: %s-context
kind: Config
preferences: {}
users:
- name: your-user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: %s
`, baseURL, cluster, cluster, cluster, cluster, cluster, pluginScript)

	err = os.WriteFile(kubeconfig, []byte(kubeconfigContent), 0600)
	if err != nil {
		return fmt.Errorf("error writing kubeconfig to file: %w", err)
	}
	fmt.Println("kubeconfig updated to use the authentication plugin for cluster", cluster)

	return nil
}