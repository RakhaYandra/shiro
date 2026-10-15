package sqlite

import (
	"path/filepath"
	"testing"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

func testSummary() scanner.Summary {
	return scanner.Summary{
		Root:           "/tmp/x",
		TotalApparent:  30,
		TotalAllocated: 40,
		AllocatedKnown: true,
		Files:          2,
		Dirs:           1,
		Tree: &scanner.DirNode{
			Path:     "/tmp/x",
			Apparent: 30,
			Children: []*scanner.DirNode{{Path: "/tmp/x/sub", Apparent: 20}},
		},
	}
}

func TestSQLite_SaveAndRead(t *testing.T) {
	db := filepath.Join(t.TempDir(), "shiro.db")
	st, err := Open(db)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	id, err := st.SaveScan(testSummary(), "test-v")
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("empty scan id")
	}
	n, err := st.EntryCount(id)
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Fatal("no entries saved")
	}
}

func TestSQLite_FKEnforced(t *testing.T) {
	db := filepath.Join(t.TempDir(), "shiro.db")
	st, err := Open(db)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if err := st.InsertEntryOrphan(); err == nil {
		t.Fatal("want FK violation, got nil")
	}
}

func TestNoop_ScanWorks(t *testing.T) {
	var s NoopStore
	if _, err := s.SaveScan(testSummary(), "v"); err != nil {
		t.Fatal(err)
	}
}

func TestSQLite_FKEnforced_MultiConn(t *testing.T) {
	db := filepath.Join(t.TempDir(), "shiro.db")
	st, err := Open(db)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	done := make(chan error, 16)
	for i := 0; i < 16; i++ {
		go func() { done <- st.InsertEntryOrphan() }()
	}
	for i := 0; i < 16; i++ {
		if err := <-done; err == nil {
			t.Fatal("orphan insert succeeded on some pooled conn: FK not enforced")
		}
	}
}

func TestSQLite_FullTree(t *testing.T) {
	db := filepath.Join(t.TempDir(), "shiro.db")
	st, err := Open(db)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	sum := testSummary()
	sum.Tree.Children = []*scanner.DirNode{
		{Path: "/tmp/x/sub", Apparent: 20, Children: []*scanner.DirNode{
			{Path: "/tmp/x/sub/leaf", Apparent: 10},
		}},
	}
	id, err := st.SaveScan(sum, "v")
	if err != nil {
		t.Fatal(err)
	}
	n, err := st.EntryCount(id)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("entries = %d, want 3 (root+sub+leaf)", n)
	}
}

func TestSQLite_ListAndGet(t *testing.T) {
	db := filepath.Join(t.TempDir(), "shiro.db")
	st, err := Open(db)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	id1, err := st.SaveScan(testSummary(), "v")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.SaveScan(testSummary(), "v"); err != nil {
		t.Fatal(err)
	}
	scans, err := st.ListScans("/tmp/x")
	if err != nil {
		t.Fatal(err)
	}
	if len(scans) != 2 || scans[0].StartedAt < scans[1].StartedAt {
		t.Fatalf("want 2 scans desc, got %+v", scans)
	}
	entries, err := st.GetEntries(id1)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("no entries returned")
	}
}
