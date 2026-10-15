CREATE TABLE IF NOT EXISTS scans (
    id TEXT PRIMARY KEY,
    root_path TEXT NOT NULL,
    started_at INTEGER NOT NULL,
    finished_at INTEGER,
    status TEXT NOT NULL,
    scanner_version TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS entries (
    scan_id TEXT NOT NULL,
    path TEXT NOT NULL,
    kind TEXT NOT NULL,
    apparent_bytes INTEGER NOT NULL,
    allocated_bytes INTEGER,
    modified_at INTEGER,
    PRIMARY KEY (scan_id, path),
    FOREIGN KEY (scan_id) REFERENCES scans(id)
);
