package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

type SqliteDriver struct {
	cfg *Config
}

func (s *SqliteDriver) ExecuteStatement(dbname string, statement string) (string, error) {
	dbPath := dbname
	if s.cfg.DBPath != "" {
		dbPath = filepath.Join(s.cfg.DBPath, dbname)
	}

	cmd := exec.Command("sqlite3", dbPath, statement, "-json")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		slog.Error("Failed to execute sqlite statement", "db", dbPath, "error", err, "stderr", stderr.String())
		return "", fmt.Errorf("sqlite error: %w. Details: %s", err, stderr.String())
	}

	return out.String(), nil
}

func (s *SqliteDriver) ExtractDDL(dbname string) (string, error) {
	dbPath := dbname
	if s.cfg.DBPath != "" {
		dbPath = filepath.Join(s.cfg.DBPath, dbname)
	}

	cmd := exec.Command("sqlite3", dbPath, ".schema")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		slog.Error("Failed to extract sqlite DDL", "db", dbPath, "error", err, "stderr", stderr.String())
		return "", fmt.Errorf("sqlite DDL error: %w. Details: %s", err, stderr.String())
	}

	return out.String(), nil
}

func (s *SqliteDriver) ListDatabases() (string, error) {
	if s.cfg.DBPath == "" {
		return "[]", nil
	}

	files, err := os.ReadDir(s.cfg.DBPath)
	if err != nil {
		return "", fmt.Errorf("failed to list sqlite databases: %w", err)
	}

	var dbs []string
	for _, f := range files {
		if !f.IsDir() && (filepath.Ext(f.Name()) == ".db" || filepath.Ext(f.Name()) == ".sqlite") {
			dbs = append(dbs, f.Name())
		}
	}

	data, _ := json.Marshal(dbs)
	return string(data), nil
}
