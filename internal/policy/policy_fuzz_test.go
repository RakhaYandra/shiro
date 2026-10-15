package policy

import "testing"

func FuzzContainment(f *testing.F) {
	for _, seed := range []string{"/a/b", "/a/../b", "/", "/a//b/", "/a/./b", "..", "../x", "/x/../../etc"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		c := canonical(s)
		if c == "" {
			t.Fatal("canonical must never be empty")
		}
		if !withinOrEqual(c, c) {
			t.Fatalf("not reflexive: %q", c)
		}
		p := Policy{Strict: []string{"/home/u/.ssh"}}
		_ = p.CheckIn(s, "/home/u/work")
	})
}
