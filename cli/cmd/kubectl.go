package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func configureKubectl(cluster string) error {
	fmt.Println("Configuring kubectl for cluster:", cluster)
	clusterFile := filepath.Join(os.Getenv("HOME"), ".kube", "cluster")
	err := os.WriteFile(clusterFile, []byte(cluster), 0600)
	if err != nil {
		fmt.Println("Error writing cluster to file:", err)
		return fmt.Errorf("error writing cluster to file: %w", err)
	}

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
		fmt.Println("Error writing plugin script to file:", err)
		return fmt.Errorf("error writing plugin script to file: %w", err)
	}
	fmt.Println("Authentication plugin script created at", pluginScript)

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
		fmt.Println("Error writing kubeconfig to file:", err)
		return fmt.Errorf("error writing kubeconfig to file: %w", err)
	}
	fmt.Println("kubeconfig updated to use the authentication plugin for cluster", cluster)

	return nil
}
