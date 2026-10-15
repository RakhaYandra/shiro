package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Walk streams entries under root without following symlinks.
// Non-fatal errors go to the second channel; both channels are closed.
func Walk(root string, opts Options) (<-chan Entry, <-chan ScanError) {
	entries := make(chan Entry, 256)
	errs := make(chan ScanError, 64)
	go func() {
		defer close(entries)
		defer close(errs)
		abs, err := filepath.Abs(root)
		if err != nil {
			errs <- ScanError{Path: root, Err: err.Error()}
			return
		}
		if resolved, err := filepath.EvalSymlinks(abs); err == nil {
			abs = resolved
		}
		excludes := make([]string, 0, len(opts.Excludes))
		for _, x := range opts.Excludes {
			if a, err := filepath.Abs(x); err == nil {
				excludes = append(excludes, a)
			}
		}
		var rootDev uint64
		if opts.OneFileSystem {
			if d, err := deviceOf(abs); err == nil {
				rootDev = d
			} else {
				errs <- ScanError{Path: abs, Err: err.Error()}
				return
			}
		}
		_ = filepath.WalkDir(abs, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				errs <- ScanError{Path: path, Err: err.Error()}
				if d != nil && d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if path != abs && isExcluded(path, excludes) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if opts.OneFileSystem {
				if dev, err := deviceOf(path); err == nil && dev != rootDev {
					if d.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}
			info, err := d.Info()
			if err != nil {
				errs <- ScanError{Path: path, Err: err.Error()}
				return nil
			}
			entries <- meta(path, info)
			return nil
		})
	}()
	return entries, errs
}

func isExcluded(path string, excludes []string) bool {
	for _, x := range excludes {
		if path == x || strings.HasPrefix(path, x+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}
