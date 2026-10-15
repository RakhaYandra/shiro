package output

import (
	"fmt"
	"math"
)

func Human(b int64) string {
	const u = 1024
	if b < u {
		return fmt.Sprintf("%d B", b)
	}
	v := float64(b)
	for _, s := range []string{"KiB", "MiB", "GiB", "TiB", "PiB"} {
		v /= u
		if v < u || s == "PiB" {
			return fmt.Sprintf("%.1f %s", v, s)
		}
	}
	return fmt.Sprintf("%d B", b)
}

func ParseSize(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	var mult int64 = 1
	num := s
	for _, suf := range []struct {
		s string
		m int64
	}{
		{"KiB", 1 << 10}, {"MiB", 1 << 20}, {"GiB", 1 << 30}, {"TiB", 1 << 40},
		{"KB", 1000}, {"MB", 1000 * 1000}, {"GB", 1000 * 1000 * 1000},
	} {
		if len(s) > len(suf.s) && s[len(s)-len(suf.s):] == suf.s {
			mult = suf.m
			num = s[:len(s)-len(suf.s)]
			break
		}
	}
	var f float64
	if _, err := fmt.Sscanf(num, "%g", &f); err != nil {
		return 0, fmt.Errorf("invalid size %q", s)
	}
	if f < 0 || math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, fmt.Errorf("invalid size %q", s)
	}
	v := f * float64(mult)
	if v > float64(math.MaxInt64) {
		return 0, fmt.Errorf("invalid size %q", s)
	}
	return int64(v), nil
}
