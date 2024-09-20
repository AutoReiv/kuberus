package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func waitForToken() (string, error) {
	fmt.Println("Waiting for token...")

	tokenFilePath := filepath.Join(os.Getenv("HOME"), ".kube", "token")
	fmt.Println("Token file path:", tokenFilePath)

	err := os.MkdirAll(filepath.Dir(tokenFilePath), 0700)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	if _, err := os.Stat(tokenFilePath); os.IsNotExist(err) {
		_, err := os.Create(tokenFilePath)
		if err != nil {
			return "", fmt.Errorf("failed to create token file: %w", err)
		}
	}

	for {
		token, err := os.ReadFile(tokenFilePath)
		if err != nil {
			fmt.Println("Failed to read token file:", err)
			return "", fmt.Errorf("failed to read token file: %w", err)
		}

		fmt.Println("Polling server for token...")
		resp, err := http.Get(fmt.Sprintf("%s/auth/token?token=%s", baseURL, string(token)))
		if err != nil {
			fmt.Println("Error making GET request:", err)
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var result map[string]string
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				fmt.Println("Error decoding JSON response:", err)
				return "", err
			}
			fmt.Println("Token retrieved successfully")
			return result["token"], nil
		}

		fmt.Println("Token not yet available, retrying in 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}

func fetchAccessibleClusters(token string) ([]string, error) {
	fmt.Println("Fetching accessible clusters...")
	req, err := http.NewRequest("GET", baseURL+"/api/clusters", nil)
	if err != nil {
		fmt.Println("Error creating new request:", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Failed to fetch clusters:", string(body))
		return nil, fmt.Errorf("failed to fetch clusters: %s", string(body))
	}

	var clusters []string
	err = json.NewDecoder(resp.Body).Decode(&clusters)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return nil, err
	}

	fmt.Println("Clusters fetched successfully:", clusters)
	return clusters, nil
}