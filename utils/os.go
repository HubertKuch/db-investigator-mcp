package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

var cacheDirOverride string

// GetCacheDir returns the path to the cache directory and ensures it exists.
func GetCacheDir() (string, error) {
	if cacheDirOverride != "" {
		if err := os.MkdirAll(cacheDirOverride, 0755); err != nil {
			return "", fmt.Errorf("failed to create override cache directory: %w", err)
		}
		return cacheDirOverride, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	cacheDir := filepath.Join(home, ".cache", "db-investigator")

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cacheDir, nil
}

// SetCacheDirOverride sets a temporary directory for cache during tests.
func SetCacheDirOverride(path string) {
	cacheDirOverride = path
}
