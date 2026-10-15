package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

func testSummary() scanner.Summary {
	return scanner.Summary{
		Root:          "/x",
		TotalApparent: 1000,
		Tree: &scanner.DirNode{
			Path: "/x", Apparent: 1000,
			Children: []*scanner.DirNode{
				{Path: "/x/b", Apparent: 700, Children: []*scanner.DirNode{
					{Path: "/x/b/leaf", Apparent: 600},
				}},
				{Path: "/x/a", Apparent: 300},
			},
		},
	}
}

func keyMsg(s string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func TestModel_Navigate(t *testing.T) {
	m := NewModel(testSummary(), nil)
	if m.Selected() != "/x" {
		t.Fatalf("selected = %q", m.Selected())
	}
	m2, _ := m.Update(keyMsg("j"))
	m = m2.(Model)
	if m.Selected() != "/x/b" {
		t.Fatalf("after j selected = %q", m.Selected())
	}
	m2, _ = m.Update(keyMsg("k"))
	m = m2.(Model)
	if m.Selected() != "/x" {
		t.Fatalf("after k selected = %q", m.Selected())
	}
}

func TestModel_ExpandCollapse(t *testing.T) {
	m := NewModel(testSummary(), nil)
	if m.VisibleRows() != 3 {
		t.Fatalf("rows = %d, want 3 (root+b+a)", m.VisibleRows())
	}
	var m2 tea.Model
	m2, _ = m.Update(keyMsg("j"))
	m = m2.(Model)
	m2, _ = m.Update(keyMsg("l"))
	m = m2.(Model)
	if m.VisibleRows() != 4 {
		t.Fatalf("rows = %d, want 4 after expanding b", m.VisibleRows())
	}
	m2, _ = m.Update(keyMsg("h"))
	m = m2.(Model)
	if m.VisibleRows() != 3 {
		t.Fatalf("rows = %d, want 3 after collapsing b", m.VisibleRows())
	}
}

func TestModel_TabSwitch(t *testing.T) {
	m := NewModel(testSummary(), nil)
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = m2.(Model)
	if m.ActiveTab() != TabFindings {
		t.Fatalf("tab = %d, want findings", m.ActiveTab())
	}
}

func TestModel_Quit(t *testing.T) {
	m := NewModel(testSummary(), nil)
	_, cmd := m.Update(keyMsg("q"))
	if cmd == nil {
		t.Fatal("q must quit")
	}
}
