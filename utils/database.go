package utils

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Crate SQLite database
func CreateDatabase(dbPath string) (*sql.DB, error) {
	// Ensure the directory exists
	db, error := sql.Open("sqlite3", dbPath)
	if error != nil {
		return nil, error
	}
	err := initDatabase(db)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func initDatabase(db *sql.DB) error {
	// Array with the table creation statements
	tableCreationQueries := []string{
		`CREATE TABLE jobs (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			source_type TEXT NOT NULL,      -- nfs, local, smb
			source_path TEXT NOT NULL,
			destination_type TEXT NOT NULL, -- s3
			destination_uri TEXT NOT NULL,
			hash_algorithm TEXT NOT NULL,   -- blake3, sha256
			enabled BOOLEAN NOT NULL DEFAULT 1,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE job_runs (
			id INTEGER PRIMARY KEY,
			job_id INTEGER NOT NULL,
			started_at DATETIME NOT NULL,
			finished_at DATETIME,
			status TEXT NOT NULL, -- running, success, failed
			files_scanned INTEGER DEFAULT 0,
			files_new INTEGER DEFAULT 0,
			files_modified INTEGER DEFAULT 0,
			files_deleted INTEGER DEFAULT 0,
			error_message TEXT,
			FOREIGN KEY (job_id) REFERENCES jobs(id)
		);`,
		`CREATE TABLE files (
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
		`CREATE TABLE file_versions (
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
		`CREATE TABLE file_state (
			file_id INTEGER PRIMARY KEY,
			current_version_id INTEGER NOT NULL,
			last_backup_run_id INTEGER NOT NULL,
			s3_key TEXT NOT NULL,
			storage_class TEXT NOT NULL, -- STANDARD, GLACIER, DEEP_ARCHIVE
			FOREIGN KEY (file_id) REFERENCES files(id)
		);`,
		`CREATE TABLE events (
			id INTEGER PRIMARY KEY,
			job_run_id INTEGER,
			level TEXT NOT NULL, -- info, warn, error
			message TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (job_run_id) REFERENCES job_runs(id)
		);`,
		// Indexes for performance
		`CREATE INDEX idx_files_job_path ON files(job_id, path);`,
		`CREATE INDEX idx_files_last_seen ON files(job_id, last_seen_run_id);`,
		`CREATE INDEX idx_files_deleted ON files(job_id, deleted_at);`,
		`CREATE INDEX idx_versions_file ON file_versions(file_id);`,
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

// Detele from files if files was deleted from source
func DeleteFileRecord(db *sql.DB, jobID int, filePath string) error {
	// UPDATE files
	// SET deleted_at = CURRENT_TIMESTAMP
	// WHERE job_id = ?
	// 	AND last_seen_run_id < ?
	// 	AND deleted_at IS NULL;
	return nil
}
