package output

import "testing"

func FuzzParseSize(f *testing.F) {
	for _, seed := range []string{"", "0", "100MiB", "2GiB", "512KiB", "1024", "1.5GiB", "bogus", "-5MiB", "99999999999999999999X"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		n, err := ParseSize(s)
		if err == nil && n < 0 {
			t.Fatalf("negative size without error: %q -> %d", s, n)
		}
		if s == "" && (err != nil || n != 0) {
			t.Fatalf("empty must be 0,nil: %d,%v", n, err)
		}
	})
}
