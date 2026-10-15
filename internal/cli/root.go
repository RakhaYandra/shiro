package cli

import "github.com/spf13/cobra"

var version = "dev"

func Version() string { return version }

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:          "shiro",
		Short:        "Know your storage. Clean with confidence.",
		Version:      version,
		SilenceUsage: true,
	}
	root.AddCommand(NewScanCmd())
	root.AddCommand(NewExplainCmd())
	root.AddCommand(NewCleanCmd())
	root.AddCommand(NewDuplicatesCmd())
	root.AddCommand(NewHistoryCmd())
	root.AddCommand(NewDiffCmd())
	root.AddCommand(NewDockerCmd())
	root.AddCommand(NewTuiCmd())
	return root
}

func Execute() int {
	if err := NewRootCmd().Execute(); err != nil {
		return ExitCode(err)
	}
	return 0
}
