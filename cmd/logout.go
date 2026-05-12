package cmd

import (
	"fmt"
	"os"

	"aspera-terminal-ui/internal/config"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Aspera system",
	Long:  `Remove saved authentication tokens.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		cfg.AccessToken = ""
		cfg.RefreshToken = ""
		cfg.ExpiresAt = ""
		cfg.RefreshExpAt = ""

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Printf("Failed to clear config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Logged out successfully.")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
