package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aspera-terminal-ui",
	Short: "Aspera Terminal UI",
	Long: `aspera-terminal-ui is a command line client for Aspera P2P/MySpace systems.
	
To enable shell completion, run:
  aspera-terminal-ui completion [bash|zsh|fish|powershell]
  
Example for bash:
  source <(aspera-terminal-ui completion bash)
  
To see help for any command:
  aspera-terminal-ui [command] --help`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
}
