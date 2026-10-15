package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RakhaYandra/shiro/internal/cleanup"
)

// RenderPlanText prints a dry-run plan with denials. Never implies execution.
func RenderPlanText(p cleanup.CleanupPlan) string {
	var b strings.Builder
	b.WriteString("Shiro · Cleanup plan (DRY-RUN, nothing removed)\n\n")
	fmt.Fprintf(&b, "Plan  %s  root %s\n\n", p.ID, p.Root)
	if len(p.Items) == 0 {
		b.WriteString("No cleanup actions planned.\n")
	}
	var total int64
	for _, it := range p.Items {
		total += it.ExpectedSize
		fmt.Fprintf(&b, "  REMOVE  %8s  [%s/%s] %s\n    %s\n",
			Human(it.ExpectedSize), it.RuleID, it.Risk, it.Path, it.Reason)
	}
	if len(p.Items) > 0 {
		fmt.Fprintf(&b, "\nEstimated reclaimable: %s (not guaranteed)\n", Human(total))
	}
	if len(p.Denied) > 0 {
		b.WriteString("\nDENIED by policy:\n")
		for _, d := range p.Denied {
			fmt.Fprintf(&b, "  %s  [%s] %s\n", d.Path, d.RuleID, d.Reason)
		}
	}
	return b.String()
}

// RenderPlanJSON emits the plan (schema_version 1).
func RenderPlanJSON(p cleanup.CleanupPlan) []byte {
	raw, _ := json.MarshalIndent(map[string]any{
		"schema_version": 1,
		"id":             p.ID,
		"root":           p.Root,
		"created_at":     p.CreatedAt,
		"dry_run":        true,
		"items":          p.Items,
		"denied":         p.Denied,
	}, "", "  ")
	return append(raw, '\n')
}
