package cli

import "testing"

func TestCleanCmd_Exists(t *testing.T) {
	root := NewRootCmd()
	for _, c := range root.Commands() {
		if c.Name() == "clean" {
			return
		}
	}
	t.Fatal("clean subcommand not registered")
}

func TestClean_RequiresDryRun(t *testing.T) {
	cmd := NewCleanCmd()
	cmd.SetArgs([]string{"/tmp"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("want error without --dry-run (no executor yet)")
	}
}
