package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/RakhaYandra/shiro/internal/cleanup"
	"github.com/RakhaYandra/shiro/internal/detector"
	"github.com/RakhaYandra/shiro/internal/output"
	"github.com/RakhaYandra/shiro/internal/policy"
	"github.com/RakhaYandra/shiro/internal/scanner"
)

var ErrNoDryRun = errors.New("refusing to clean without --dry-run")

var ErrPolicyDenied = errors.New("cleanup denied by policy or validation")

var ErrPartialFail = errors.New("cleanup partially failed")

func NewCleanCmd() *cobra.Command {
	var dryRun, asJSON, execute, yes bool
	var category, minSizeStr, planFile string
	cmd := &cobra.Command{
		Use:   "clean [PATH]",
		Short: "Preview or execute reviewed cleanup",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if planFile != "" {
				return runExecute(cmd, planFile, execute, yes, asJSON)
			}
			if !dryRun {
				return fmt.Errorf("%w: use --dry-run or --plan FILE --execute", ErrNoDryRun)
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
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview only")
	cmd.Flags().StringVar(&planFile, "plan", "", "plan file from dry-run --json")
	cmd.Flags().BoolVar(&execute, "execute", false, "execute the plan file")
	cmd.Flags().BoolVar(&yes, "yes", false, "skip interactive confirmation (policy still enforced)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "machine-readable output")
	cmd.Flags().StringVar(&category, "category", "", "filter by category")
	cmd.Flags().StringVar(&minSizeStr, "min-size", "", "filter e.g. 100MiB")
	return cmd
}

func runExecute(cmd *cobra.Command, planFile string, execute, yes, asJSON bool) error {
	if !execute {
		return fmt.Errorf("%w: --plan requires --execute", ErrNoDryRun)
	}
	plan, err := cleanup.LoadPlan(planFile)
	if err != nil {
		return err
	}
	if len(plan.Items) == 0 {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Nothing to execute: plan has no items.")
		return nil
	}
	if !yes {
		if !isTTY(cmd.InOrStdin()) {
			return fmt.Errorf("%w: refusing without a terminal; use --yes", ErrNoDryRun)
		}
		_, _ = fmt.Fprintf(cmd.OutOrStdout(),
			"Execute plan %s (%d items)? Type 'yes' to confirm: ", plan.ID, len(plan.Items))
		line, _ := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
		if strings.TrimSpace(strings.ToLower(line)) != "yes" {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
			return nil
		}
	}
	home, _ := os.UserHomeDir()
	rep := cleanup.Execute(policy.Default(home), plan)
	if asJSON {
		raw, _ := json.MarshalIndent(rep, "", "  ")
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(raw))
	} else {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), output.RenderReportText(rep))
	}
	switch {
	case rep.Failed > 0:
		return ErrPartialFail
	case rep.Success == 0:
		return ErrPolicyDenied
	default:
		return nil
	}
}

func isTTY(f interface{}) bool {
	st, ok := f.(*os.File)
	if !ok {
		return false
	}
	fi, err := st.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
