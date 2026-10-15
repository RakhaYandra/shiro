package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/RakhaYandra/shiro/internal/detector"
	"github.com/RakhaYandra/shiro/internal/output"
	"github.com/RakhaYandra/shiro/internal/scanner"
)

type Tab int

const (
	TabTree Tab = iota
	TabFindings
	TabPlan
)

type row struct {
	node  *scanner.DirNode
	depth int
}

type Model struct {
	sum      scanner.Summary
	findings []detector.Finding
	plan     planView
	cursor   int
	expanded map[string]bool
	sortSize bool
	tab      Tab
	rows     []row
	width    int
	height   int
}

func NewModel(sum scanner.Summary, findings []detector.Finding) Model {
	m := Model{
		sum:      sum,
		findings: findings,
		expanded: map[string]bool{},
		sortSize: true,
	}
	if sum.Tree != nil {
		m.expanded[sum.Tree.Path] = true
	}
	m.refresh()
	return m
}

func (m *Model) refresh() {
	m.rows = nil
	if m.sum.Tree == nil {
		return
	}
	var walk func(n *scanner.DirNode, depth int)
	walk = func(n *scanner.DirNode, depth int) {
		m.rows = append(m.rows, row{n, depth})
		if !m.expanded[n.Path] {
			return
		}
		kids := append([]*scanner.DirNode(nil), n.Children...)
		if m.sortSize {
			sort.Slice(kids, func(i, j int) bool { return kids[i].Apparent > kids[j].Apparent })
		} else {
			sort.Slice(kids, func(i, j int) bool { return kids[i].Path < kids[j].Path })
		}
		for _, k := range kids {
			walk(k, depth+1)
		}
	}
	walk(m.sum.Tree, 0)
	if m.cursor >= len(m.rows) {
		m.cursor = max(len(m.rows)-1, 0)
	}
}

func (m Model) Selected() string {
	if len(m.rows) == 0 {
		return ""
	}
	return m.rows[m.cursor].node.Path
}

func (m Model) VisibleRows() int { return len(m.rows) }

func (m Model) ActiveTab() Tab { return m.tab }

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			if m.tab == TabTree && m.cursor < len(m.rows)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.tab == TabTree && m.cursor > 0 {
				m.cursor--
			}
		case "l", "right", "enter":
			if m.tab == TabTree && len(m.rows) > 0 {
				m.expanded[m.rows[m.cursor].node.Path] = true
				m.refresh()
			}
		case "h", "left":
			if m.tab == TabTree && len(m.rows) > 0 {
				path := m.rows[m.cursor].node.Path
				if path != m.sum.Root {
					delete(m.expanded, path)
					m.refresh()
				}
			}
		case "s":
			m.sortSize = !m.sortSize
			m.refresh()
		case "tab":
			m.tab = (m.tab + 1) % 3
		}
	}
	return m, nil
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString(m.tabs())
	b.WriteString("\n")
	switch m.tab {
	case TabTree:
		b.WriteString(m.treeView())
	case TabFindings:
		b.WriteString(m.findingsView())
	case TabPlan:
		b.WriteString(m.plan.View())
	}
	b.WriteString("\nj/k move · h/l collapse/expand · s sort · tab switch · q quit\n")
	return b.String()
}

func (m Model) tabs() string {
	names := []string{"tree", "findings", "plan"}
	var parts []string
	for i, n := range names {
		if Tab(i) == m.tab {
			parts = append(parts, "["+n+"]")
		} else {
			parts = append(parts, " "+n+" ")
		}
	}
	return strings.Join(parts, " ")
}

func (m Model) treeView() string {
	var b strings.Builder
	for i, r := range m.rows {
		mark := " "
		if i == m.cursor {
			mark = ">"
		}
		fmt.Fprintf(&b, "%s %s%s  %s\n", mark,
			strings.Repeat("  ", r.depth), shortBase(r.node.Path),
			output.Human(r.node.Apparent))
	}
	return b.String()
}

func shortBase(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 && i < len(path)-1 {
		return path[i+1:]
	}
	return path
}
