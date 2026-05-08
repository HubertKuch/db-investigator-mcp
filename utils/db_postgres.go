package utils

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

type PostgresDriver struct {
	cfg *Config
}

func (p *PostgresDriver) ExecuteStatement(dbname string, statement string) (string, error) {
	cmd := exec.Command("psql",
		"-h", p.cfg.Host,
		"-p", p.cfg.Port,
		"-U", p.cfg.User,
		"-d", dbname,
		"-c", statement,
		"-q", "",
		"-t", "",
	)

	cmd.Env = append(os.Environ(), "PGPASSWORD="+p.cfg.Password)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		slog.Error("Failed to execute statement", "error", err, "stderr", stderr.String())
		return "", fmt.Errorf("error executing query: %w. Details: %s", err, stderr.String())
	}

	return out.String(), nil
}

func (p *PostgresDriver) ExtractDDL(dbname string) (string, error) {
	cmd := exec.Command("pg_dump",
		"-h", p.cfg.Host,
		"-p", p.cfg.Port,
		"-U", p.cfg.User,
		"-s", dbname,
	)

	cmd.Env = append(os.Environ(), "PGPASSWORD="+p.cfg.Password)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		slog.Error("Failed to extract DDL", "error", err, "stderr", stderr.String())
		return "", fmt.Errorf("error extracting DDL: %w. Details: %s", err, stderr.String())
	}

	return out.String(), nil
}

func (p *PostgresDriver) ListDatabases() (string, error) {
	query := "SELECT json_agg(datname) FROM pg_database WHERE datistemplate = false AND datname != 'postgres';"
	return p.ExecuteStatement("postgres", query)
}
