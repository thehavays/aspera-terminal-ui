package cmd

import (
	"fmt"
	"os"

	"aspera-terminal-ui/internal/config"
	"aspera-terminal-ui/internal/tui"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List transfers",
	Long:  `List received or shared transfers using the interactive TUI.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil || cfg.AccessToken == "" {
			fmt.Println("Not logged in. Please run 'atui login' first.")
			os.Exit(1)
		}

		received, _ := cmd.Flags().GetBool("received")
		shared, _ := cmd.Flags().GetBool("shared")

		lt := tui.Received
		if shared && !received {
			lt = tui.Shared
		}

		err = tui.RunListTUI(lt)
		if err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("received", "r", false, "List received transfers (default)")
	listCmd.Flags().BoolP("shared", "s", false, "List shared transfers")
}
