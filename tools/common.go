package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"praca-db-tools-mcp/utils"
)

// saveToCache saves the given content to a file in the global cache directory.
func saveToCache(filename string, content string) (string, error) {
	cacheDir, err := utils.GetCacheDir()
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(cacheDir, filename)
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write cache file: %w", err)
	}

	return fullPath, nil
}

// readFromCache reads the content of a file from the global cache directory.
func readFromCache(filename string) (string, error) {
	cacheDir, err := utils.GetCacheDir()
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(cacheDir, filename)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
