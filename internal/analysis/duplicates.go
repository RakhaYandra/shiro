package analysis

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

type DupOptions struct {
	MinSize  int64
	Excludes []string
}

type DupGroup struct {
	Hash  string
	Size  int64
	Files []string
}

type DupReport struct {
	Root        string
	Groups      []DupGroup
	Reclaimable int64
	FilesHashed int64
	Errors      []scanner.ScanError
	Partial     bool
}

type dupRef struct {
	path string
	dev  uint64
	ino  uint64
}

func FindDuplicates(root string, opts DupOptions) (DupReport, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return DupReport{}, err
	}
	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		abs = resolved
	}
	rep := DupReport{Root: abs}
	entries, errch := scanner.Walk(abs, scanner.Options{Excludes: opts.Excludes})
	bySize := map[int64][]dupRef{}
	var errs []scanner.ScanError
	done := make(chan struct{})
	go func() {
		for e := range errch {
			errs = append(errs, e)
		}
		close(done)
	}()
	seen := map[[2]uint64]bool{}
	for e := range entries {
		if e.Kind != scanner.KindFile || e.ApparentBytes < opts.MinSize {
			continue
		}
		key := [2]uint64{e.DeviceID, e.Inode}
		if e.DeviceID+e.Inode != 0 {
			if seen[key] {
				continue
			}
			seen[key] = true
		}
		bySize[e.ApparentBytes] = append(bySize[e.ApparentBytes], dupRef{e.Path, e.DeviceID, e.Inode})
	}
	<-done
	rep.Errors = errs
	rep.Partial = len(errs) > 0
	for size, refs := range bySize {
		if len(refs) < 2 {
			continue
		}
		for _, g := range hashGroups(size, refs, &rep) {
			if len(g.Files) < 2 {
				continue
			}
			rep.Groups = append(rep.Groups, g)
			rep.Reclaimable += int64(len(g.Files)-1) * size
		}
	}
	return rep, nil
}

const prefixLen = 64 << 10

func hashGroups(size int64, refs []dupRef, rep *DupReport) []DupGroup {
	byPrefix := map[string][]dupRef{}
	for _, r := range refs {
		h, ok := hashPrefix(r.path)
		if !ok {
			rep.Errors = append(rep.Errors, scanner.ScanError{Path: r.path, Err: "unreadable during hash"})
			rep.Partial = true
			continue
		}
		rep.FilesHashed++
		byPrefix[h] = append(byPrefix[h], r)
	}
	var out []DupGroup
	for _, sub := range byPrefix {
		if len(sub) < 2 {
			continue
		}
		byFull := map[string][]string{}
		for _, r := range sub {
			h, ok := hashFull(r.path, size)
			if !ok {
				rep.Errors = append(rep.Errors, scanner.ScanError{Path: r.path, Err: "unreadable during hash"})
				rep.Partial = true
				continue
			}
			rep.FilesHashed++
			byFull[h] = append(byFull[h], r.path)
		}
		for h, paths := range byFull {
			if len(paths) < 2 {
				continue
			}
			out = append(out, DupGroup{Hash: h, Size: size, Files: paths})
		}
	}
	return out
}

func hashPrefix(path string) (string, bool) {
	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.CopyN(h, f, prefixLen); err != nil && err != io.EOF {
		return "", false
	}
	return hex.EncodeToString(h.Sum(nil)), true
}

func hashFull(path string, size int64) (string, bool) {
	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil || fi.Size() != size {
		return "", false
	}
	h := sha256.New()
	n, err := io.Copy(h, f)
	if err != nil || n != size {
		return "", false
	}
	return hex.EncodeToString(h.Sum(nil)), true
}
