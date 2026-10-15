package detector

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

type Finding struct {
	RuleID          string
	Path            string
	Bytes           int64
	Category        string
	Risk            string
	Confidence      string
	Reason          string
	SuggestedAction string
}

type Rule struct {
	ID       string
	Category string
	Basename string
	Parent   string
	Markers  []string
	Describe string
}

type Registry struct {
	rules []Rule
	stat  func(string) (os.FileInfo, error)
}

func NewRegistry() *Registry {
	return &Registry{stat: os.Stat}
}

func (r *Registry) Register(rule Rule) {
	r.rules = append(r.rules, rule)
}

func (r *Registry) Detect(sum scanner.Summary) []Finding {
	var out []Finding
	if sum.Tree == nil {
		return nil
	}
	var walk func(n *scanner.DirNode)
	walk = func(n *scanner.DirNode) {
		for _, rule := range r.rules {
			if filepath.Base(n.Path) != rule.Basename {
				continue
			}
			if rule.Parent != "" && filepath.Base(filepath.Dir(n.Path)) != rule.Parent {
				continue
			}
			marker := ""
			for _, m := range rule.Markers {
				if _, err := r.stat(filepath.Join(n.Path, m)); err == nil {
					marker = m + " (inside)"
					break
				}
				if _, err := r.stat(filepath.Join(filepath.Dir(n.Path), m)); err == nil {
					marker = m
					break
				}
			}
			conf := "low"
			reason := fmt.Sprintf("%s matched %s (%s); no marker file found",
				n.Path, rule.Basename, rule.Describe)
			if marker != "" {
				conf = "high"
				reason = fmt.Sprintf("%s matched %s (%s); marker %s present",
					n.Path, rule.Basename, rule.Describe, marker)
			} else if len(rule.Markers) == 0 {
				conf = "medium"
				reason = fmt.Sprintf("%s matched %s (%s)",
					n.Path, rule.Basename, rule.Describe)
			}
			out = append(out, Finding{
				RuleID:          rule.ID,
				Path:            n.Path,
				Bytes:           n.Apparent,
				Category:        rule.Category,
				Risk:            "review",
				Confidence:      conf,
				Reason:          reason,
				SuggestedAction: "review contents before cleanup (planner phase)",
			})
		}
		for _, k := range n.Children {
			walk(k)
		}
	}
	walk(sum.Tree)
	return out
}
