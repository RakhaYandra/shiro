package cli

import "github.com/spf13/cobra"

var version = "dev"

// Version returns the build version (overridden via ldflags).
func Version() string { return version }

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:          "shiro",
		Short:        "Know your storage. Clean with confidence.",
		Version:      version,
		SilenceUsage: true,
	}
	root.AddCommand(NewScanCmd())
	return root
}

// Execute runs the root command, mapping errors to exit codes.
// 0 ok, 1 general error, 2 invalid CLI usage, 3 partial scan.
func Execute() int {
	if err := NewRootCmd().Execute(); err != nil {
		return ExitCode(err)
	}
	return 0
}
