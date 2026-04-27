package sa

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

type Cache struct {
	db *sql.DB
}

func OpenCache(path string) (*Cache, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open cache: %w", err)
	}

	// Create tables if they don't exist
	schema := `
	CREATE TABLE IF NOT EXISTS release_cache (
		mbid TEXT PRIMARY KEY,
		last_checked TIMESTAMP,
		etag TEXT,
		data BLOB
	);
	CREATE TABLE IF NOT EXISTS file_cache (
		path TEXT PRIMARY KEY,
		mtime TIMESTAMP,
		size INTEGER,
		fingerprint TEXT,
		duration INTEGER,
		acoustid TEXT
	);
	`
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &Cache{db: db}, nil
}

func (c *Cache) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *Cache) GetFingerprint(path string, mtime time.Time, size int64) (fp string, dur int, aid string, err error) {
	row := c.db.QueryRow("SELECT fingerprint, duration, acoustid FROM file_cache WHERE path = ? AND mtime = ? AND size = ?", path, mtime, size)
	err = row.Scan(&fp, &dur, &aid)
	return
}

func (c *Cache) SaveFingerprint(path string, mtime time.Time, size int64, fp string, dur int, aid string) error {
	_, err := c.db.Exec("INSERT OR REPLACE INTO file_cache (path, mtime, size, fingerprint, duration, acoustid) VALUES (?, ?, ?, ?, ?, ?)",
		path, mtime, size, fp, dur, aid)
	return err
}

func (c *Cache) GetRelease(mbid string) (data []byte, lastChecked time.Time, etag string, err error) {
	row := c.db.QueryRow("SELECT data, last_checked, etag FROM release_cache WHERE mbid = ?", mbid)
	err = row.Scan(&data, &lastChecked, &etag)
	return
}

func (c *Cache) SaveRelease(mbid string, data []byte, etag string) error {
	_, err := c.db.Exec("INSERT OR REPLACE INTO release_cache (mbid, last_checked, etag, data) VALUES (?, ?, ?, ?)",
		mbid, time.Now(), etag, data)
	return err
}
