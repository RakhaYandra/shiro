package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/RakhaYandra/shiro/internal/analysis"
	"github.com/RakhaYandra/shiro/internal/output"
)

func NewDuplicatesCmd() *cobra.Command {
	var asJSON bool
	var minSizeStr string
	var excludes []string
	cmd := &cobra.Command{
		Use:   "duplicates PATH",
		Short: "Find duplicate file content (read-only)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			minSize, err := output.ParseSize(minSizeStr)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrUsage, err)
			}
			rep, err := analysis.FindDuplicates(args[0], analysis.DupOptions{
				MinSize:  minSize,
				Excludes: excludes,
			})
			if err != nil {
				return err
			}
			if asJSON {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(output.RenderDupJSON(rep)))
			} else {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), output.RenderDupText(rep))
			}
			if rep.Partial {
				return ErrPartial
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "machine-readable output")
	cmd.Flags().StringVar(&minSizeStr, "min-size", "1MiB", "only hash files at/above this size")
	cmd.Flags().StringArrayVar(&excludes, "exclude", nil, "exclude path (repeatable)")
	return cmd
}
