package cli

import "testing"

func TestHistoryDiffCmds_Exist(t *testing.T) {
	root := NewRootCmd()
	want := map[string]bool{"history": false, "diff": false}
	for _, c := range root.Commands() {
		if _, ok := want[c.Name()]; ok {
			want[c.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Fatalf("%s subcommand not registered", name)
		}
	}
}

func TestDiff_BadIDs(t *testing.T) {
	cmd := NewDiffCmd()
	cmd.SetArgs([]string{"--from", "nope", "--to", "nada", "--db", "/does-not-exist-shiro/x.db"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("want error for unknown scan IDs")
	}
}
