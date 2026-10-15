package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/RakhaYandra/shiro/internal/docker"
	"github.com/RakhaYandra/shiro/internal/output"
)

func NewDockerCmd() *cobra.Command {
	var host string
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "docker",
		Short: "Docker storage inventory (read-only)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sock := docker.DefaultSocket()
			if host != "" {
				sock = host
			}
			df, err := docker.NewClient(sock).SystemDF()
			if err != nil {
				return err
			}
			rep := docker.Summarize(df)
			if asJSON {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(output.RenderDockerJSON(rep)))
			} else {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), output.RenderDockerText(rep))
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "machine-readable output")
	cmd.Flags().StringVar(&host, "host", "", "docker socket path override")
	return cmd
}
