package utils

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cespare/xxhash/v2"
	"github.com/zeebo/blake3"
)

// Checks if a file or folder exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		// File exists
		return true
	}
	if os.IsNotExist(err) {
		// File does not exist
		return false
	}
	// Some other error (permission etc.)
	return false
}

// Checks if a folder exists
func FolderExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// Calculate a file checksum passing the algorithm as parameter
// default algoritmh will be xxhash
func ChecksumFile(path string, algorithm string) (string, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	switch algorithm {
	case "blake3":
		h := blake3.New()
		if _, err := io.Copy(h, file); err != nil {
			return "", err
		}
		return hex.EncodeToString(h.Sum(nil)), nil
	case "xxhash":
		h := xxhash.New()
		if _, err := io.Copy(h, file); err != nil {
			return "", err
		}
		return fmt.Sprintf("%016x", h.Sum64()), nil

	default:
		return "", fmt.Errorf("Unsupported hash algorithm: %s", algorithm)
	}
}

// Walk directory
func WalkDirectory(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return GetFileInfo(path, info)
		}
		return nil
	})
}

func GetFileInfo(path string, info os.FileInfo) error {
	// Example callback function that just prints the file path and size
	hash, err := ChecksumFile(path, "xxhash")
	if err != nil {
		return err
	}
	fmt.Printf("File: %s, Size: %d bytes, Hash: %s\n", path, info.Size(), hash)
	return nil
}
