package analysis

import (
	"sort"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

// Category is one storage bucket with traceable evidence.
type Category struct {
	Category string
	Bytes    int64
	Files    int64
	RuleID   string
	Reason   string
}

var extCategory = map[string]string{
	// media
	".mp4": "media", ".mkv": "media", ".avi": "media", ".mov": "media",
	".mp3": "media", ".flac": "media", ".ogg": "media", ".wav": "media",
	".jpg": "media", ".jpeg": "media", ".png": "media", ".gif": "media",
	".webp": "media", ".svg": "media",
	// archives
	".zip": "archives", ".tar": "archives", ".gz": "archives",
	".xz": "archives", ".7z": "archives", ".rar": "archives",
	".deb": "archives", ".rpm": "archives", ".pkg": "archives",
	// source
	".go": "source", ".rs": "source", ".py": "source", ".js": "source",
	".ts": "source", ".tsx": "source", ".jsx": "source", ".php": "source",
	".java": "source", ".c": "source", ".h": "source", ".cpp": "source",
	".rb": "source", ".lua": "source", ".sh": "source", ".sql": "source",
	// docs
	".md": "documents", ".txt": "documents", ".pdf": "documents",
	".docx": "documents", ".odt": "documents", ".tex": "documents",
	// data
	".db": "data", ".sqlite": "data", ".sqlite3": "data", ".json": "data",
	".csv": "data", ".parquet": "data", ".yaml": "data", ".yml": "data",
	".toml": "data",
	// binaries
	".o": "binaries", ".a": "binaries", ".so": "binaries",
	".exe": "binaries", ".elf": "binaries", ".wasm": "binaries",
}

// Categorize folds extension totals into named buckets, largest first.
// Unknown extensions and extensionless files go to "other".
func Categorize(sum scanner.Summary) []Category {
	buckets := map[string]*Category{}
	get := func(name, rule, reason string) *Category {
		c, ok := buckets[name]
		if !ok {
			c = &Category{Category: name, RuleID: rule, Reason: reason}
			buckets[name] = c
		}
		return c
	}
	for ext, bytes := range sum.ExtTotals {
		count := sum.ExtCounts[ext]
		name, ok := extCategory[ext]
		if !ok {
			c := get("other", "CAT-000", "extension not in known map (or none)")
			c.Bytes += bytes
			c.Files += count
			continue
		}
		c := get(name, "CAT-EXT", "extension "+displayExt(ext)+" mapped to "+name)
		c.Bytes += bytes
		c.Files += count
	}
	out := make([]Category, 0, len(buckets))
	for _, c := range buckets {
		out = append(out, *c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Bytes > out[j].Bytes })
	return out
}

func displayExt(ext string) string {
	if ext == "" {
		return "(none)"
	}
	return ext
}
