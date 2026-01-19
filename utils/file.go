package utils

import "os"

// fileExists checks if a file or folder exists
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
