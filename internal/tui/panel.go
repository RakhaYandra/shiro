package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RakhaYandra/shiro/internal/output"
	"github.com/RakhaYandra/shiro/internal/policy"
	"github.com/RakhaYandra/shiro/internal/scanner"
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

func (m Model) nodeDetail(path string) string {
	var b strings.Builder
	node := findNode(m.sum.Tree, path)
	if node == nil {
		return "(node not found)\n"
	}
	fmt.Fprintf(&b, "%s  %s (%d files)\n", node.Path, output.Human(node.Apparent), node.Files)
	for _, f := range m.findings {
		if !within(path, f.Path) {
			continue
		}
		home, _ := os.UserHomeDir()
		dec := policy.Decide(policy.Default(home), f, m.sum.Root)
		verdict := "offered (needs confirmation)"
		if !dec.Allowed {
			verdict = "blocked: " + dec.Reason
		}
		fmt.Fprintf(&b, "  [%s/%s/%s] %s\n    %s\n    policy: %s\n",
			f.RuleID, f.Confidence, f.Risk, f.Path, f.Reason, verdict)
	}
	return b.String()
}

func findNode(n *scanner.DirNode, path string) *scanner.DirNode {
	if n == nil {
		return nil
	}
	if n.Path == path {
		return n
	}
	for _, k := range n.Children {
		if found := findNode(k, path); found != nil {
			return found
		}
	}
	return nil
}

func within(dir, path string) bool {
	return path == dir || strings.HasPrefix(path, dir+string(filepath.Separator))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
