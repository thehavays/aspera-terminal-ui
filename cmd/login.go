package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"aspera-terminal-ui/internal/api"
	"aspera-terminal-ui/internal/config"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Aspera system",
	Long:  `Authenticate and save token for further requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		reader := bufio.NewReader(os.Stdin)

		if cfg.Endpoint == "" {
			fmt.Print("Endpoint (e.g., https://api.aspera.com): ")
			cfg.Endpoint, _ = reader.ReadString('\n')
			cfg.Endpoint = strings.TrimSpace(cfg.Endpoint)
		}

		if cfg.ClientID == "" {
			fmt.Print("Client ID: ")
			cfg.ClientID, _ = reader.ReadString('\n')
			cfg.ClientID = strings.TrimSpace(cfg.ClientID)
		}

		if cfg.ClientSecret == "" {
			fmt.Print("Client Secret: ")
			cfg.ClientSecret, _ = reader.ReadString('\n')
			cfg.ClientSecret = strings.TrimSpace(cfg.ClientSecret)
		}

		fmt.Print("Username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Print("Password: ")
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Printf("\nError reading password: %v\n", err)
			os.Exit(1)
		}
		password := string(bytePassword)
		fmt.Println()

		fmt.Println("Authenticating...")

		tResp, err := api.FetchToken(cfg.Endpoint, username, password, cfg.ClientID, cfg.ClientSecret)
		if err != nil {
			fmt.Printf("Login failed: %v\n", err)
			os.Exit(1)
		}

		cfg.Username = username
		cfg.Password = password
		cfg.AccessToken = tResp.AccessToken
		cfg.RefreshToken = tResp.RefreshToken
		cfg.ExpiresAt = time.Now().Add(time.Duration(tResp.ExpiresIn) * time.Second).Format(time.RFC3339)
		cfg.RefreshExpAt = time.Now().Add(time.Duration(tResp.RefreshTokenExpiresIn) * time.Second).Format(time.RFC3339)

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Printf("Failed to save config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Login successful! Token saved.")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
