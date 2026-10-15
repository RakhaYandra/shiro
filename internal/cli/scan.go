package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/RakhaYandra/shiro/internal/output"
	"github.com/RakhaYandra/shiro/internal/scanner"
)

// ErrPartial signals a scan that completed with non-fatal errors (exit 3).
var ErrPartial = errors.New("partial scan")

// ErrUsage signals invalid CLI usage (exit 2).
var ErrUsage = errors.New("invalid usage")

// ExitCode maps errors to the documented CLI contract:
// 0 ok, 1 general error, 2 invalid CLI usage, 3 partial scan.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	if errors.Is(err, ErrPartial) {
		return 3
	}
	if errors.Is(err, ErrUsage) {
		return 2
	}
	msg := err.Error()
	for _, s := range []string{"unknown flag", "accepts ", "requires ", "invalid argument"} {
		if strings.Contains(msg, s) {
			return 2
		}
	}
	return 1
}

type ScanFlags struct {
	JSON       bool
	Depth      int
	Excludes   []string
	OneFileSys bool
	MinSize    string
	NoStore    bool
	DBPath     string
}

func NewScanCmd() *cobra.Command {
	var f ScanFlags
	cmd := &cobra.Command{
		Use:   "scan PATH",
		Short: "Scan storage (read-only)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			minSize, err := output.ParseSize(f.MinSize)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrUsage, err)
			}
			sum, err := scanner.Aggregate(args[0], scanner.Options{
				Excludes:      f.Excludes,
				OneFileSystem: f.OneFileSys,
			})
			if err != nil {
				return err
			}
			if f.JSON {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(output.RenderJSON(sum, version)))
			} else {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), output.RenderText(sum, output.DisplayOpts{
					Depth:   f.Depth,
					MinSize: minSize,
				}))
			}
			if sum.Partial {
				return ErrPartial
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&f.JSON, "json", false, "machine-readable output")
	cmd.Flags().IntVar(&f.Depth, "depth", 3, "display depth")
	cmd.Flags().StringArrayVar(&f.Excludes, "exclude", nil, "exclude path (repeatable)")
	cmd.Flags().BoolVar(&f.OneFileSys, "one-file-system", false, "do not cross filesystem boundaries")
	cmd.Flags().StringVar(&f.MinSize, "min-size", "", "filter e.g. 100MiB")
	cmd.Flags().BoolVar(&f.NoStore, "no-store", false, "skip sqlite inventory")
	cmd.Flags().StringVar(&f.DBPath, "db", "", "sqlite path override")
	return cmd
}
