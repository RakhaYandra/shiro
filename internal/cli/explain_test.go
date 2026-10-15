package cli

import "testing"

func TestExplainCmd_Exists(t *testing.T) {
	root := NewRootCmd()
	for _, c := range root.Commands() {
		if c.Name() == "explain" {
			return
		}
	}
	t.Fatal("explain subcommand not registered")
}

func TestExplain_BadPath(t *testing.T) {
	cmd := NewExplainCmd()
	cmd.SetArgs([]string{"/does-not-exist-shiro"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("want error for missing path")
	}
}
