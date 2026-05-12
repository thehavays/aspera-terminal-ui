package cmd

import (
	"fmt"
	"os"
	"strings"

	"aspera-terminal-ui/internal/api"
	"aspera-terminal-ui/internal/config"
	"aspera-terminal-ui/internal/tui"

	"golang.org/x/term"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) > max {
		return string(runes[:max-3]) + "..."
	}
	return s
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List transfers",
	Long:  `List received or shared transfers.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil || cfg.AccessToken == "" {
			fmt.Println("Not logged in. Please run 'atui login' first.")
			os.Exit(1)
		}

		received, _ := cmd.Flags().GetBool("received")
		shared, _ := cmd.Flags().GetBool("shared")
		interactive, _ := cmd.Flags().GetBool("interactive")

		// If no flags provided or interactive explicitly requested, use TUI
		if interactive || (!received && !shared) {
			lt := tui.Received
			if shared && !received {
				lt = tui.Shared
			}
			err := tui.RunListTUI(lt)
			if err != nil {
				fmt.Println("Error running TUI:", err)
				os.Exit(1)
			}
			return
		}

		filter := api.ListFilter{
			PageIndex: 1,
			PageSize:  50,
		}

		// Dinamik terminal genişliği
		termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil || termWidth <= 0 {
			termWidth = 120 // Varsayılan genişlik
		}

		if received {
			fmt.Println("\n=== RECEIVED TRANSFERS ===")
			resp, err := api.QueryReceivedRequests(cfg.Endpoint, cfg.AccessToken, filter)
			if err != nil {
				fmt.Printf("Error fetching received transfers: %v\n", err)
				os.Exit(1)
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Request ID", "Title", "Sender", "Send Time", "Status", "DL'd"})

			t.SetStyle(table.StyleColoredDark)
			t.Style().Options.DrawBorder = true
			t.Style().Options.SeparateRows = false

			senderMax := 25
			if termWidth > 150 {
				senderMax = 35
			}
			titleMax := termWidth - (63 + senderMax + 2)
			if titleMax < 20 {
				titleMax = 20
			}

			for _, item := range resp.DataList {
				dlStr := "No"
				if item.Downloaded {
					dlStr = "Yes"
				}

				title := truncate(strings.TrimSpace(item.Title), titleMax)
				sender := truncate(strings.TrimSpace(item.UserEmail), senderMax)

				t.AppendRow(table.Row{
					item.RequestId,
					title,
					sender,
					item.SendTime,
					item.Status,
					dlStr,
				})
			}
			t.Render()
			fmt.Printf("Total Received: %d\n\n", resp.Count)
		}

		if shared {
			fmt.Println("\n=== SHARED TRANSFERS ===")
			resp, err := api.QuerySharedRequests(cfg.Endpoint, cfg.AccessToken, filter)
			if err != nil {
				fmt.Printf("Error fetching shared transfers: %v\n", err)
				os.Exit(1)
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Request ID", "Subject", "Share Time", "Status", "System"})

			t.SetStyle(table.StyleColoredDark)
			t.Style().Options.DrawBorder = true
			t.Style().Options.SeparateRows = false

			subjectMax := termWidth - 66
			if subjectMax < 20 {
				subjectMax = 20
			}

			for _, item := range resp.DataList {
				subject := truncate(strings.TrimSpace(item.Subject), subjectMax)

				t.AppendRow(table.Row{
					item.RequestId,
					subject,
					item.ShareTime,
					item.Status,
					item.System,
				})
			}
			t.Render()
			fmt.Printf("Total Shared: %d\n\n", resp.Count)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("received", "r", false, "List received transfers")
	listCmd.Flags().BoolP("shared", "s", false, "List shared transfers")
	listCmd.Flags().BoolP("interactive", "i", false, "Force interactive TUI mode")
}
