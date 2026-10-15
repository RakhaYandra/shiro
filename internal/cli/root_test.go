package cli

import "testing"

func TestScanCmd_Exists(t *testing.T) {
	root := NewRootCmd()
	found := false
	for _, c := range root.Commands() {
		if c.Name() == "scan" {
			found = true
		}
	}
	if !found {
		t.Fatal("scan subcommand not registered")
	}
}

func TestRoot_VersionSet(t *testing.T) {
	root := NewRootCmd()
	if root.Version == "" {
		t.Fatal("version must be set")
	}
}
