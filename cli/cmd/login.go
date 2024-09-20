package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the application and configure kubectl",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting login process...")

		// Debug statement to check if the baseURL is being set
		fmt.Println("Debug: baseURL in login command is set to", baseURL)

		if baseURL == "" {
			fmt.Println("Error: Base URL is not set. Please provide the --url flag.")
			os.Exit(1)
		}

		loginURL := fmt.Sprintf("%s/login", baseURL)
		fmt.Println("Opening browser to URL:", loginURL)
		err := openBrowser(loginURL)
		if err != nil {
			fmt.Println("Error opening browser:", err)
			os.Exit(1)
		}

		token, err := waitForToken()
		if err != nil {
			fmt.Println("Error retrieving token:", err)
			os.Exit(1)
		}

		clusters, err := fetchAccessibleClusters(token)
		if err != nil {
			fmt.Println("Error fetching accessible clusters:", err)
			os.Exit(1)
		}

		for _, cluster := range clusters {
			fmt.Println("Configuring kubectl for cluster:", cluster)
			err := configureKubectl(cluster)
			if err != nil {
				fmt.Println("Error configuring kubectl for cluster", cluster, ":", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(LoginCmd)
}