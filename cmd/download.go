package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"path/filepath"

	"aspera-terminal-ui/internal/api"
	"aspera-terminal-ui/internal/aspera"
	"aspera-terminal-ui/internal/config"
	"aspera-terminal-ui/internal/tui"
	"github.com/spf13/cobra"
)

var (
	destDir string
	allFiles bool
)

var downloadCmd = &cobra.Command{
	Use:   "download [REQUEST_ID]",
	Short: "Download files for a received or shared request using Aspera",
	Args:  cobra.MaximumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		
		cfg, err := config.LoadConfig()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		
		token, err := api.EnsureValidToken()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		
		resp, err := api.QueryReceivedRequests(cfg.Endpoint, token, api.ListFilter{PageIndex: 1, PageSize: 20})
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		
		var suggestions []string
		for _, req := range resp.DataList {
			suggestions = append(suggestions, fmt.Sprintf("%s\t%s", req.RequestId, req.Title))
		}
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		var reqID string
		var selectedFiles []string

		if len(args) == 0 {
			// Interactive Request Selection
			var err error
			reqID, err = tui.SelectRequest()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			if reqID == "" {
				fmt.Println("No request selected.")
				return
			}
		} else {
			reqID = args[0]
		}

		// Interactive File Selection (if not --all)
		if !allFiles {
			var err error
			selectedFiles, err = tui.SelectFiles(reqID)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}

		targetDir := destDir
		if targetDir == "" {
			var err error
			targetDir, err = os.Getwd()
			if err != nil {
				fmt.Println("Error getting current directory:", err)
				os.Exit(1)
			}
		}

		absTargetDir, err := filepath.Abs(targetDir)
		if err != nil {
			fmt.Println("Error resolving target directory:", err)
			os.Exit(1)
		}

		// Ensure directory exists
		if err := os.MkdirAll(absTargetDir, 0755); err != nil {
			fmt.Println("Error creating target directory:", err)
			os.Exit(1)
		}

		fmt.Printf("Fetching download instructions for request: %s\n", reqID)
		
		dlInfo, err := api.GetDownloadCommand(reqID, selectedFiles)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		ascpPath, keyPath, err := aspera.GetAscpPaths()
		if err != nil {
			fmt.Println("Failed to initialize Aspera binaries:", err)
			os.Exit(1)
		}

		ascpArgsTmpl := dlInfo.Results.AscpArgs
		cmdStr := strings.ReplaceAll(ascpArgsTmpl, "{{ascp.exe path}}", ascpPath)
		cmdStr = strings.ReplaceAll(cmdStr, "{ascp.exe path}", ascpPath)
		cmdStr = strings.ReplaceAll(cmdStr, "{{openssh file path}}", keyPath)
		cmdStr = strings.ReplaceAll(cmdStr, "{openssh file path}", keyPath)
		
		cmdStr = strings.ReplaceAll(cmdStr, "{download dir}", fmt.Sprintf("\"%s\"", absTargetDir))
		cmdStr = strings.ReplaceAll(cmdStr, "{{download dir}}", fmt.Sprintf("\"%s\"", absTargetDir))

		fmt.Printf("Starting Aspera download to: %s\n", absTargetDir)
		
		c := exec.Command("sh", "-c", cmdStr)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		
		if err := c.Run(); err != nil {
			fmt.Println("\nDownload failed:", err)
			os.Exit(1)
		}
		
		fmt.Println("\nDownload completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&destDir, "dest", "d", "", "Destination directory for downloaded files (default: current directory)")
	downloadCmd.Flags().BoolVarP(&allFiles, "all", "a", false, "Download all files without asking")
}
