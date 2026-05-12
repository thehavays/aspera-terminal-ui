package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "atui",
	Short: "Aspera Terminal UI",
	Long: `atui is a command line client for Aspera P2P/MySpace systems.
	
To enable shell completion, run:
  atui completion [bash|zsh|fish|powershell]
  
Example for bash:
  source <(atui completion bash)
  
To see help for any command:
  atui [command] --help`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
}
