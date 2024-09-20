package main

import (
	"fmt"
	"os"

	"rbac/cli/cmd"

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
	rootCmd.AddCommand(cmd.LoginCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}