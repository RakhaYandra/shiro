package cli

import (
	"errors"

	"github.com/spf13/cobra"
)

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
			return errors.New("scan not implemented yet")
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
