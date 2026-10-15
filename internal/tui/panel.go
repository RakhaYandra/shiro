package tui

import (
	"fmt"
	"strings"
)

func (m Model) findingsView() string {
	if len(m.findings) == 0 {
		return "(no developer artifacts detected)\n"
	}
	var b strings.Builder
	for _, f := range m.findings {
		fmt.Fprintf(&b, "  [%s/%s] %s\n    %s\n", f.RuleID, f.Confidence, f.Path, f.Reason)
	}
	return b.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
