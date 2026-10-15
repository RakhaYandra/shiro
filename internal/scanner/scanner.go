package scanner

type FileKind string

const (
	KindFile      FileKind = "file"
	KindDirectory FileKind = "directory"
	KindSymlink   FileKind = "symlink"
	KindOther     FileKind = "other"
)

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

type ScanError struct {
	Path string
	Err  string
}

type Options struct {
	Excludes      []string
	OneFileSystem bool
}
