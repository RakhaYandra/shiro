package cli

import "testing"

func TestExitCode_PolicyAndPartial(t *testing.T) {
	if got := ExitCode(ErrPolicyDenied); got != 4 {
		t.Fatalf("got %d, want 4", got)
	}
	if got := ExitCode(ErrPartialFail); got != 5 {
		t.Fatalf("got %d, want 5", got)
	}
}

func TestExecute_MissingPlan(t *testing.T) {
	cmd := NewCleanCmd()
	cmd.SetArgs([]string{"--plan", "/does-not-exist-shiro/plan.json", "--execute", "--yes"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("want error for missing plan")
	}
}

func TestExecute_NeedsConfirm(t *testing.T) {
	cmd := NewCleanCmd()
	cmd.SetArgs([]string{"--plan", "/does-not-exist-shiro/plan.json", "--execute"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("want error: missing plan must fail before prompt")
	}
}
