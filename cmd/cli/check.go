package cli

import (
	"403unlocker-cli/internal/check"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks if the DNS SNI-Proxy can bypass 403 error for a specific domain",
	Long: `Checks if the DNS SNI-Proxy can bypass 403 error for a specific domain

Examples:
    403unlocker check https://pkg.go.dev`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !check.DomainValidator(args[0]) {
			cmd.Printf("Error: '%s' is not a valid domain format\n\n", args[0])
			return cmd.Help()
		}
		return check.CheckWithDNS(args[0])
	},
	Aliases: []string{"c"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
