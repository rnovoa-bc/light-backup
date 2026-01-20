package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// OPen SQLite database
func OpenDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Crate SQLite database
func CreateDatabase(dbPath string, initialize *bool) error {
	// Ensure the directory exists
	log.Println("Creating database at:", dbPath)
	if !FileExists(dbPath) {
		err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("Error creating directories: %v", err)
		}
	}
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	// Initialize schema if requested
	if *initialize {
		if err := DeleteSchema(db); err != nil {
			return fmt.Errorf("Error deleting schema: %v", err)
		}
	}
	err = CreateSchema(db)
	if err != nil {
		return fmt.Errorf("Error creating schema: %v", err)
	}
	return nil
}

func CreateSchema(db *sql.DB) error {
	// Array with the table creation statements
	log.Println("Creating database schema...")

	tableCreationQueries := []string{
		`CREATE TABLE IF NOT EXISTS destinations (
			id INTEGER PRIMARY KEY,	
			region TEXT NOT NULL,
			bucket TEXT NOT NULL,
			access_key TEXT NOT NULL,
			secret_key TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS jobs (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			source_path TEXT NOT NULL,
			source_type TEXT NOT NULL, -- local, nfs, etc.
			source_options TEXT,       -- JSON encoded options
			destination_id INTEGER NOT NULL, -- s3 destination reference
			hash_algorithm TEXT NOT NULL,   -- blake3, sha256
			max_duration_seconds INTEGER NOT NULL,
			enabled BOOLEAN NOT NULL DEFAULT 1,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			foreign KEY (destination_id) REFERENCES destinations(id)
		);`,
		`CREATE TABLE IF NOT EXISTS job_runs (
			id INTEGER PRIMARY KEY,
			job_id INTEGER NOT NULL,
			started_at_seconds INTEGER NOT NULL,
			finished_at_seconds INTEGER,
			status TEXT NOT NULL, -- running, success, failed
			files_scanned INTEGER DEFAULT 0,
			files_new INTEGER DEFAULT 0,
			files_modified INTEGER DEFAULT 0,
			files_deleted INTEGER DEFAULT 0,
			error_message TEXT,
			FOREIGN KEY (job_id) REFERENCES jobs(id)
		);`,
		`CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY,
			job_id INTEGER NOT NULL,
			path TEXT NOT NULL,
			size INTEGER NOT NULL,
			mtime DATETIME,
			last_seen_run_id INTEGER NOT NULL,
			deleted_at DATETIME,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			UNIQUE(job_id, path),
			FOREIGN KEY (job_id) REFERENCES jobs(id)
		);`,
		`CREATE TABLE IF NOT EXISTS file_versions (
			id INTEGER PRIMARY KEY,
			file_id INTEGER NOT NULL,
			size INTEGER NOT NULL,
			hash TEXT NOT NULL,
			hash_algorithm TEXT NOT NULL,
			first_seen_run_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			UNIQUE(file_id, hash),
			FOREIGN KEY (file_id) REFERENCES files(id)
		);`,
		`CREATE TABLE IF NOT EXISTS file_state (
			file_id INTEGER PRIMARY KEY,
			current_version_id INTEGER NOT NULL,
			last_backup_run_id INTEGER NOT NULL,
			s3_key TEXT NOT NULL,
			storage_class TEXT NOT NULL, -- STANDARD, GLACIER, DEEP_ARCHIVE
			FOREIGN KEY (file_id) REFERENCES files(id)
		);`,

		// Indexes for performance
		`CREATE INDEX IF NOT EXISTS idx_files_job_path ON files(job_id, path);`,
		`CREATE INDEX IF NOT EXISTS idx_files_last_seen ON files(job_id, last_seen_run_id);`,
		`CREATE INDEX IF NOT EXISTS idx_files_deleted ON files(job_id, deleted_at);`,
		`CREATE INDEX IF NOT EXISTS idx_versions_file ON file_versions(file_id);`,
	}

	// Execute each table creation statement
	for _, query := range tableCreationQueries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

// Deslete existing tables if needed
func DeleteSchema(db *sql.DB) error {
	log.Println("Deleting database schema...")
	tables := []string{
		"file_state",
		"file_versions",
		"files",
		"job_runs",
		"jobs",
		"destinations",
	}
	for _, table := range tables {
		_, err := db.Exec("DROP TABLE IF EXISTS " + table + ";")
		if err != nil {
			return err
		}
	}
	return nil
}

// Detele from files if files was deleted from source
func DeleteFileRecord(db *sql.DB, jobID int, filePath string) error {
	// UPDATE files
	// SET deleted_at = CURRENT_TIMESTAMP
	// WHERE job_id = ?
	// 	AND last_seen_run_id < ?
	// 	AND deleted_at IS NULL;
	return nil
}

// Adds a new destination to the database
func AddDestination(db *sql.DB, region, bucket, accessKey, secretKey string) (int64, error) {
	result, err := db.Exec(
		`INSERT INTO destinations (region, bucket, access_key, secret_key, created_at, updated_at)
		 VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);`,
		region, bucket, accessKey, secretKey,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
