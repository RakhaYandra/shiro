package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/RakhaYandra/shiro/internal/analysis"
	"github.com/RakhaYandra/shiro/internal/output"
	"github.com/RakhaYandra/shiro/internal/scanner"
)

func NewExplainCmd() *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "explain [PATH]",
		Short: "Explain storage usage (read-only)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) == 1 {
				root = args[0]
			}
			sum, err := scanner.Aggregate(root, scanner.Options{})
			if err != nil {
				return err
			}
			if asJSON {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(),
					string(output.RenderExplainJSON(sum, analysis.Reclaimable(sum), version)))
			} else {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), output.RenderExplainText(sum, version))
			}
			if sum.Partial {
				return ErrPartial
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "machine-readable output")
	return cmd
}
