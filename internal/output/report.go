package output

import (
	"fmt"
	"strings"

	"github.com/RakhaYandra/shiro/internal/cleanup"
)

func RenderReportText(rep cleanup.Report) string {
	var b strings.Builder
	b.WriteString("Shiro · Cleanup report\n\n")
	fmt.Fprintf(&b, "Plan %s: %d success, %d skipped, %d denied, %d failed\n",
		rep.PlanID, rep.Success, rep.Skipped, rep.Denied, rep.Failed)
	fmt.Fprintf(&b, "Reclaimed: %s (verified removals only)\n", Human(rep.BytesReclaimed))
	for _, r := range rep.Results {
		fmt.Fprintf(&b, "  %-8s %s\n    %s\n", r.Status, r.Path, r.Reason)
	}
	return b.String()
}
