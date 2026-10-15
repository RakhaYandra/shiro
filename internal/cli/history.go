package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/RakhaYandra/shiro/internal/analysis"
	"github.com/RakhaYandra/shiro/internal/output"
	"github.com/RakhaYandra/shiro/internal/store/sqlite"
)

func openStore(dbPath string) (*sqlite.Store, string, error) {
	if dbPath == "" {
		dbPath = defaultDBPath()
	}
	if dbPath == "" {
		return nil, "", fmt.Errorf("no home directory for default store; use --db")
	}
	st, err := sqlite.Open(dbPath)
	return st, dbPath, err
}

func canonicalRoot(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		return resolved
	}
	return abs
}

func NewHistoryCmd() *cobra.Command {
	var dbPath string
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "history [PATH]",
		Short: "List recorded scans (read-only)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) == 1 {
				root = args[0]
			}
			root = canonicalRoot(root)
			st, _, err := openStore(dbPath)
			if err != nil {
				return err
			}
			defer st.Close()
			scans, err := st.ListScans(root)
			if err != nil {
				return err
			}
			if jsonOut {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(output.RenderHistoryJSON(root, scans)))
			} else {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), output.RenderHistoryText(root, scans))
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "machine-readable output")
	cmd.Flags().StringVar(&dbPath, "db", "", "sqlite path override")
	return cmd
}

func NewDiffCmd() *cobra.Command {
	var from, to, dbPath string
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "diff [PATH]",
		Short: "Compare two scans (read-only)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, _, err := openStore(dbPath)
			if err != nil {
				return err
			}
			defer st.Close()
			if from == "" || to == "" {
				root := "."
				if len(args) == 1 {
					root = args[0]
				}
				root = canonicalRoot(root)
				scans, err := st.ListScans(root)
				if err != nil {
					return err
				}
				if len(scans) < 2 {
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Need at least two scans for this path. Run `shiro scan` twice.")
					return nil
				}
				from, to = scans[1].ID, scans[0].ID
			}
			fromRoot, err := st.ScanRoot(from)
			if err != nil {
				return fmt.Errorf("unknown --from scan: %w", err)
			}
			toRoot, err := st.ScanRoot(to)
			if err != nil {
				return fmt.Errorf("unknown --to scan: %w", err)
			}
			if fromRoot != toRoot {
				return fmt.Errorf("%w: scans cover different roots (%s vs %s)", ErrUsage, fromRoot, toRoot)
			}
			oldEntries, err := st.GetEntries(from)
			if err != nil {
				return err
			}
			newEntries, err := st.GetEntries(to)
			if err != nil {
				return err
			}
			d := analysis.Diff(oldEntries, newEntries)
			if jsonOut {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(output.RenderDiffJSON(from, to, d)))
			} else {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), output.RenderDiffText(from, to, d))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "older scan ID")
	cmd.Flags().StringVar(&to, "to", "", "newer scan ID")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "machine-readable output")
	cmd.Flags().StringVar(&dbPath, "db", "", "sqlite path override")
	return cmd
}
