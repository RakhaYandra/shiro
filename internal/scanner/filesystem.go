package scanner

import (
	"os"
	"syscall"
)

// meta extracts portable metadata from an lstat result.
func meta(path string, fi os.FileInfo) Entry {
	e := Entry{Path: path, ModTimeUnix: fi.ModTime().Unix()}
	mode := fi.Mode()
	switch {
	case mode.IsRegular():
		e.Kind = KindFile
		e.ApparentBytes = fi.Size()
	case mode.IsDir():
		e.Kind = KindDirectory
	case (mode & os.ModeSymlink) != 0:
		e.Kind = KindSymlink
	default:
		e.Kind = KindOther
	}
	if st, ok := fi.Sys().(*syscall.Stat_t); ok && st != nil {
		e.DeviceID = uint64(st.Dev)
		e.Inode = st.Ino
		if e.Kind == KindFile {
			e.AllocatedBytes = st.Blocks * 512
			e.AllocatedKnown = true
		}
	}
	return e
}

// deviceOf returns the device ID for boundary checks (stubbed in tests if needed).
func deviceOf(path string) (uint64, error) {
	var st syscall.Stat_t
	if err := syscall.Stat(path, &st); err != nil {
		return 0, err
	}
	return uint64(st.Dev), nil
}
