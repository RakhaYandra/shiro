package cli

import "testing"

func TestDuplicatesCmd_Exists(t *testing.T) {
	root := NewRootCmd()
	for _, c := range root.Commands() {
		if c.Name() == "duplicates" {
			return
		}
	}
	t.Fatal("duplicates subcommand not registered")
}

func TestDuplicates_BadSize(t *testing.T) {
	cmd := NewDuplicatesCmd()
	cmd.SetArgs([]string{"/tmp", "--min-size", "bogus"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("want error for invalid --min-size")
	}
}
