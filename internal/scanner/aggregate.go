package scanner

import (
	"container/heap"
	"path/filepath"
	"sort"
)

// FileRef is one large-file record for ranking.
type FileRef struct {
	Path          string
	ApparentBytes int64
}

// topHeap is a bounded min-heap of largest files.
type topHeap []FileRef

func (h topHeap) Len() int            { return len(h) }
func (h topHeap) Less(i, j int) bool  { return h[i].ApparentBytes < h[j].ApparentBytes }
func (h topHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *topHeap) Push(x interface{}) { *h = append(*h, x.(FileRef)) }
func (h *topHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// maxTopFiles bounds ranking memory regardless of tree size.
const maxTopFiles = 100

// DirNode is one directory subtree with rolled-up sizes.
type DirNode struct {
	Path      string
	Apparent  int64
	Allocated int64
	Children  []*DirNode
	Files     int64
	Dirs      int64
}

// Summary is the full scan result.
type Summary struct {
	Root           string
	TotalApparent  int64
	TotalAllocated int64
	AllocatedKnown bool
	Files          int64
	Dirs           int64
	Symlinks       int64
	Others         int64
	Errors         []ScanError
	Partial        bool
	Tree           *DirNode
	TopFiles       []FileRef
}

// Aggregate walks root and rolls sizes bottom-up.
// Hardlinks (same dev+ino) count once. Symlinks count as entries, size 0.
func Aggregate(root string, opts Options) (Summary, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return Summary{}, err
	}
	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		abs = resolved
	}
	entries, errch := Walk(abs, opts)
	nodes := map[string]*DirNode{abs: {Path: abs}}
	seen := map[[2]uint64]bool{}
	var sum Summary
	sum.Root = abs
	var errs []ScanError
	errDone := make(chan struct{})
	go func() {
		for e := range errch {
			errs = append(errs, e)
		}
		close(errDone)
	}()
	var tops topHeap
	heap.Init(&tops)
	pushTop := func(f FileRef) {
		if tops.Len() < maxTopFiles {
			heap.Push(&tops, f)
			return
		}
		if f.ApparentBytes > tops[0].ApparentBytes {
			heap.Pop(&tops)
			heap.Push(&tops, f)
		}
	}
	for e := range entries {
		key := [2]uint64{e.DeviceID, e.Inode}
		if e.Kind == KindFile && e.DeviceID+e.Inode != 0 {
			if seen[key] {
				continue
			}
			seen[key] = true
		}
		switch e.Kind {
		case KindFile:
			sum.Files++
			sum.TotalApparent += e.ApparentBytes
			if e.AllocatedKnown {
				sum.TotalAllocated += e.AllocatedBytes
				sum.AllocatedKnown = true
			}
			pushTop(FileRef{Path: e.Path, ApparentBytes: e.ApparentBytes})
			addToParents(nodes, abs, dirOf(e.Path, abs), e.ApparentBytes, e.AllocatedBytes, e.AllocatedKnown)
		case KindDirectory:
			sum.Dirs++
			if _, ok := nodes[e.Path]; !ok {
				nodes[e.Path] = &DirNode{Path: e.Path}
			}
			nodes[e.Path].Dirs++
		case KindSymlink:
			sum.Symlinks++
		default:
			sum.Others++
		}
	}
	<-errDone
	sum.Errors = errs
	sum.Partial = len(errs) > 0
	linkChildren(nodes, abs)
	sum.Tree = nodes[abs]
	sum.TopFiles = make([]FileRef, tops.Len())
	for i := len(sum.TopFiles) - 1; i >= 0; i-- {
		sum.TopFiles[i] = heap.Pop(&tops).(FileRef)
	}
	return sum, nil
}

func dirOf(path, root string) string {
	d := filepath.Dir(path)
	if d == "." || d == "" {
		return root
	}
	return d
}

func addToParents(nodes map[string]*DirNode, root, dir string, apparent, allocated int64, known bool) {
	for {
		n, ok := nodes[dir]
		if !ok {
			n = &DirNode{Path: dir}
			nodes[dir] = n
		}
		n.Apparent += apparent
		n.Files++
		if known {
			n.Allocated += allocated
		}
		if dir == root {
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}

func linkChildren(nodes map[string]*DirNode, root string) {
	for path, n := range nodes {
		if path == root {
			continue
		}
		p := filepath.Dir(path)
		if parent, ok := nodes[p]; ok {
			parent.Children = append(parent.Children, n)
		}
	}
	for _, n := range nodes {
		sort.Slice(n.Children, func(i, j int) bool {
			return n.Children[i].Apparent > n.Children[j].Apparent
		})
	}
}
