package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"aspera-terminal-ui/internal/api"
	"aspera-terminal-ui/internal/aspera"
	"github.com/spf13/cobra"
)

var (
	shareTo      []string
	shareCc      []string
	shareBcc     []string
	shareSubject string
	shareBody    string
)

var shareCmd = &cobra.Command{
	Use:   "share [FILES...]",
	Short: "Share files via Aspera",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if shareSubject == "" {
			fmt.Println("Error: subject is required")
			os.Exit(1)
		}
		if len(shareTo) == 0 {
			fmt.Println("Error: at least one recipient (--to) is required")
			os.Exit(1)
		}

		fmt.Println("Creating share request...")
		resp, err := api.CreateShareRequest(shareSubject, shareBody, shareTo, shareCc, shareBcc, args)
		if err != nil {
			fmt.Printf("Error creating share request: %v\n", err)
			os.Exit(1)
		}

		ascpPath, keyPath, err := aspera.GetAscpPaths()
		if err != nil {
			fmt.Printf("Error getting Aspera paths: %v\n", err)
			os.Exit(1)
		}

		cmdStr := resp.Results.AscpArgs
		if cmdStr == "" {
			fmt.Println("Error: API returned empty Aspera command arguments (ascpargs)")
			os.Exit(1)
		}
		cmdStr = strings.ReplaceAll(cmdStr, "{{ascp.exe path}}", ascpPath)
		cmdStr = strings.ReplaceAll(cmdStr, "{ascp.exe path}", ascpPath)
		cmdStr = strings.ReplaceAll(cmdStr, "{{openssh file path}}", keyPath)
		cmdStr = strings.ReplaceAll(cmdStr, "{openssh file path}", keyPath)

		for _, f := range args {
			absPath, _ := filepath.Abs(f)
			basename := filepath.Base(f)
			// Replace {file dir}/basename with absolute path
			cmdStr = strings.ReplaceAll(cmdStr, "{{file dir}}/"+basename, absPath)
			cmdStr = strings.ReplaceAll(cmdStr, "{file dir}/"+basename, absPath)
		}

		fmt.Printf("Starting Aspera upload for Request ID: %s\n", resp.Results.RequestId)
		
		c := exec.Command("sh", "-c", cmdStr)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			fmt.Printf("\nUpload failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\nUpload completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)

	shareCmd.Flags().StringSliceVarP(&shareTo, "to", "t", []string{}, "Recipient email addresses (comma separated)")
	shareCmd.Flags().StringSliceVar(&shareCc, "cc", []string{}, "CC email addresses (comma separated)")
	shareCmd.Flags().StringSliceVar(&shareBcc, "bcc", []string{}, "BCC email addresses (comma separated)")
	shareCmd.Flags().StringVarP(&shareSubject, "subject", "s", "", "Subject of the share request")
	shareCmd.Flags().StringVarP(&shareBody, "body", "b", "", "Body message of the share request")

	shareCmd.MarkFlagRequired("to")
	shareCmd.MarkFlagRequired("subject")
}
