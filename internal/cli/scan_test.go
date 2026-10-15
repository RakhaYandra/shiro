package cli

import "testing"

func TestExitCode_Partial(t *testing.T) {
	if got := ExitCode(ErrPartial); got != 3 {
		t.Fatalf("got %d, want 3", got)
	}
	if got := ExitCode(nil); got != 0 {
		t.Fatalf("got %d, want 0", got)
	}
}

func TestScan_InvalidMinSize(t *testing.T) {
	cmd := NewScanCmd()
	cmd.SetArgs([]string{"/tmp", "--min-size", "bogus"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("want error for invalid --min-size")
	}
	if got := ExitCode(err); got != 2 {
		t.Fatalf("exit = %d, want 2", got)
	}
}
