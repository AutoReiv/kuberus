package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var baseURL string

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "A CLI tool for interacting with the RBAC app",
	Long:  `A CLI tool for logging in and configuring kubectl to interact with the RBAC app.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please use a subcommand. For example, 'mycli login --url http://localhost:3000'")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&baseURL, "url", "u", "", "Base URL of the RBAC app")
	rootCmd.MarkPersistentFlagRequired("url")

	// Debug statement to check if the baseURL is being set
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		fmt.Println("Debug: baseURL is set to", baseURL)
	}
}

func Execute() error {
	return rootCmd.Execute()
}