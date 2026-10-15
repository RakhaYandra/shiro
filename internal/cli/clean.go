package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/RakhaYandra/shiro/internal/cleanup"
	"github.com/RakhaYandra/shiro/internal/detector"
	"github.com/RakhaYandra/shiro/internal/output"
	"github.com/RakhaYandra/shiro/internal/policy"
	"github.com/RakhaYandra/shiro/internal/scanner"
)

var ErrNoDryRun = errors.New("refusing to clean without --dry-run")

func NewCleanCmd() *cobra.Command {
	var dryRun, asJSON bool
	var category, minSizeStr string
	cmd := &cobra.Command{
		Use:   "clean [PATH]",
		Short: "Preview cleanup (dry-run only, never deletes)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !dryRun {
				return fmt.Errorf("%w: executor not implemented", ErrNoDryRun)
			}
			root := "."
			if len(args) == 1 {
				root = args[0]
			}
			minSize, err := output.ParseSize(minSizeStr)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrUsage, err)
			}
			sum, err := scanner.Aggregate(root, scanner.Options{})
			if err != nil {
				return err
			}
			home, _ := os.UserHomeDir()
			plan, err := cleanup.Plan(policy.Default(home),
				detector.DefaultRegistry().Detect(sum), sum.Root,
				cleanup.Filter{Category: category, MinSize: minSize})
			if err != nil {
				return err
			}
			if asJSON {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(output.RenderPlanJSON(plan)))
			} else {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), output.RenderPlanText(plan))
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview only (required)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "machine-readable output")
	cmd.Flags().StringVar(&category, "category", "", "filter by category")
	cmd.Flags().StringVar(&minSizeStr, "min-size", "", "filter e.g. 100MiB")
	return cmd
}
