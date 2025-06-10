package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var Version = "dev" // <- Injected by ldflags

var rootCmd = &cobra.Command{
	Use:     "403unlocker",
	Short:   "403Unlocker-CLI is a versatile command-line tool designed to bypass 403 restrictions effectively",
	Long:    `403Unlocker-CLI is a versatile command-line tool designed to bypass 403 restrictions effectively`,
	Version: Version, // Enables --version flag
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
