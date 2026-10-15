package scanner

// FileKind classifies an inventory entry.
type FileKind string

const (
	KindFile      FileKind = "file"
	KindDirectory FileKind = "directory"
	KindSymlink   FileKind = "symlink"
	KindOther     FileKind = "other"
)

// Entry is one filesystem inventory record.
type Entry struct {
	Path           string
	Kind           FileKind
	ApparentBytes  int64
	AllocatedBytes int64
	AllocatedKnown bool
	ModTimeUnix    int64
	DeviceID       uint64
	Inode          uint64
}

// ScanError records a non-fatal per-path failure.
type ScanError struct {
	Path string
	Err  string
}

// Options controls traversal scope.
type Options struct {
	Excludes      []string
	OneFileSystem bool
}
