package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/RakhaYandra/shiro/internal/cleanup"
	"github.com/RakhaYandra/shiro/internal/detector"
	"github.com/RakhaYandra/shiro/internal/output"
	"github.com/RakhaYandra/shiro/internal/policy"
	"github.com/RakhaYandra/shiro/internal/scanner"
	"github.com/RakhaYandra/shiro/internal/tui"
)

func NewTuiCmd() *cobra.Command {
	var minSizeStr string
	cmd := &cobra.Command{
		Use:   "tui [PATH]",
		Short: "Interactive storage browser (read-only preview)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !isTerminal(os.Stdin) || !isTerminal(os.Stdout) {
				return fmt.Errorf("tui requires a terminal; use scan/explain/clean for pipes")
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
			findings := detector.DefaultRegistry().Detect(sum)
			home, _ := os.UserHomeDir()
			plan, err := cleanup.Plan(policy.Default(home), findings, sum.Root,
				cleanup.Filter{MinSize: minSize})
			if err != nil {
				return err
			}
			m := tui.NewModel(sum, findings)
			m.SetPlan(plan)
			_, err = tea.NewProgram(m).Run()
			return err
		},
	}
	cmd.Flags().StringVar(&minSizeStr, "min-size", "", "filter plan preview e.g. 100MiB")
	return cmd
}

func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
