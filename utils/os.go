package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("nie udało się pobrać katalogu domowego: %w", err)
	}

	cacheDir := filepath.Join(home, ".cache", "db-investigator")

	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return "", fmt.Errorf("nie udało się stworzyć katalogu cache: %w", err)
	}

	return cacheDir, nil
}

func EnsureCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("nie udało się pobrać katalogu domowego: %w", err)
	}

	cacheDir := filepath.Join(home, ".cache", "db-investigator")

	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return "", fmt.Errorf("nie udało się zainicjalizować katalogu cache: %w", err)
	}

	return cacheDir, nil
}
