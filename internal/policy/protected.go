package policy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Policy struct {
	Strict []string
	Roots  []string
}

func Default(home string) Policy {
	strict := []string{".ssh", ".gnupg", ".password-store"}
	roots := []string{"Documents", "Projects", ".config", ".local/share", ".local/state"}
	p := Policy{}
	for _, d := range strict {
		p.Strict = append(p.Strict, filepath.Join(home, d))
	}
	for _, d := range roots {
		p.Roots = append(p.Roots, filepath.Join(home, d))
	}
	return p
}

func canonical(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		return resolved
	}
	return abs
}

func withinOrEqual(path, base string) bool {
	return path == base || strings.HasPrefix(path, base+string(os.PathSeparator))
}

func (p Policy) Check(target string) error {
	canon := canonical(target)
	for _, s := range p.Strict {
		if withinOrEqual(canon, canonical(s)) {
			return fmt.Errorf("protected: %s is inside %s", target, s)
		}
	}
	for _, r := range p.Roots {
		if canon == canonical(r) {
			return fmt.Errorf("protected: %s is a protected root", target)
		}
	}
	return nil
}

func (p Policy) CheckIn(target, root string) error {
	if root == "" {
		return p.Check(target)
	}
	canon, canonRoot := canonical(target), canonical(root)
	if !withinOrEqual(canon, canonRoot) {
		return fmt.Errorf("out of scope: %s is outside scan root %s", target, root)
	}
	return p.Check(target)
}
