package cli

import (
	"403unlocker-cli/internal/dns"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	timeoutBestDns int
	checkfirst     bool
)

var bestdnsCmd = &cobra.Command{
	Use:   "bestdns",
	Short: "Finds the fastest DNS SNI-Proxy for downloading a specific URL",
	Long: `Finds the fastest DNS SNI-Proxy for downloading a specific URL

Examples:
	403unlocker bestdns --timeout 15 https://packages.gitlab.com/gitlab/gitlab-ce/packages/el/7/gitlab-ce-16.8.0-ce.0.el7.x86_64.rpm/download.rpm`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !dns.URLValidator(args[0]) {
			fmt.Println("Error: Invalid URL")
			return cmd.Help()
		}
		return dns.CheckWithURL(args[0], checkfirst, timeoutBestDns)
	},
	Aliases: []string{"dns"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	bestdnsCmd.PersistentFlags().IntVarP(&timeoutBestDns, "timeout", "t", 10, "Sets timeout")
	bestdnsCmd.PersistentFlags().BoolVarP(&checkfirst, "check", "c", false, "first check with GET to route '/' of domain ")

	rootCmd.AddCommand(bestdnsCmd)
}
