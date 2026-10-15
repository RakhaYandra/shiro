package cli

import "testing"

func TestTuiCmd_Exists(t *testing.T) {
	root := NewRootCmd()
	for _, c := range root.Commands() {
		if c.Name() == "tui" {
			return
		}
	}
	t.Fatal("tui subcommand not registered")
}

func TestTui_RefusesNonTTY(t *testing.T) {
	cmd := NewTuiCmd()
	cmd.SetArgs([]string{"/tmp"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("want error when not attached to a terminal")
	}
}
