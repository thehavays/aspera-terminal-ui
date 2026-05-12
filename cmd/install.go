package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"aspera-terminal-ui/internal/aspera"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install embedded Aspera Client",
	Run: func(cmd *cobra.Command, args []string) {
		ascpPath, keyPath, err := aspera.GetAscpPaths()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Ascp path: %s\nKey path: %s\n", ascpPath, keyPath)

		setupCompletion()
	},
}

func setupCompletion() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	shells := []struct {
		name string
		file string
		cmd  string
	}{
		{"bash", ".bashrc", "source <(atui completion bash)"},
		{"zsh", ".zshrc", "source <(atui completion zsh)"},
	}

	for _, s := range shells {
		path := filepath.Join(home, s.file)
		if _, err := os.Stat(path); err != nil {
			continue
		}

		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		if strings.Contains(string(content), s.cmd) {
			fmt.Printf("Shell completion already set up for %s\n", s.name)
			continue
		}

		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			continue
		}

		fmt.Fprintf(f, "\n# Aspera Terminal UI (atui) completion\n%s\n", s.cmd)
		f.Close()
		fmt.Printf("Shell completion added to %s\n", s.file)
	}
	fmt.Println("Please restart your terminal or run 'source ~/.bashrc' (or ~/.zshrc) to enable completion.")
}

func init() {
	rootCmd.AddCommand(installCmd)
}
