package cli

import (
	"403unlocker-cli/internal/docker"

	"github.com/spf13/cobra"
)

var (
	timeoutFastDocker int
)

var fastdockerCmd = &cobra.Command{
	Use:   "fastdocker",
	Short: "Finds the fastest docker registries for a specific docker image",
	Long: `Examples:
    403unlocker fastdocker --timeout 15 gitlab/gitlab-ce:17.0.0-ce.0`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if docker.DockerImageValidator(args[0]) {
			return docker.CheckWithDockerImage(args[0], timeoutFastDocker)
		} else {
			cmd.Printf("Error: '%s' is not a valid docker image format\n\n", args[0])
			return cmd.Help()
		}
	},
	Aliases: []string{"docker"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	fastdockerCmd.PersistentFlags().IntVarP(&timeoutFastDocker, "timeout", "t", 10, "Sets timeout")

	rootCmd.AddCommand(fastdockerCmd)

}
