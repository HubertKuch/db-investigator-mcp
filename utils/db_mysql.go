package utils

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
)

type MysqlDriver struct {
	cfg *Config
}

func (m *MysqlDriver) ExecuteStatement(dbname string, statement string) (string, error) {
	cmd := exec.Command("mysql",
		"-h", m.cfg.Host,
		"-P", m.cfg.Port,
		"-u", m.cfg.User,
		"-p"+m.cfg.Password,
		"-D", dbname,
		"-e", statement,
		"--batch",
		"--raw",
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		slog.Error("Failed to execute mysql statement", "db", dbname, "error", err, "stderr", stderr.String())
		return "", fmt.Errorf("mysql error: %w. Details: %s", err, stderr.String())
	}

	return out.String(), nil
}

func (m *MysqlDriver) ExtractDDL(dbname string) (string, error) {
	// mysqldump -h host -P port -u user -ppassword --no-data dbname
	cmd := exec.Command("mysqldump",
		"-h", m.cfg.Host,
		"-P", m.cfg.Port,
		"-u", m.cfg.User,
		"-p"+m.cfg.Password,
		"--no-data",
		dbname,
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		slog.Error("Failed to extract mysql DDL", "db", dbname, "error", err, "stderr", stderr.String())
		return "", fmt.Errorf("mysql DDL error: %w. Details: %s", err, stderr.String())
	}

	return out.String(), nil
}

func (m *MysqlDriver) ListDatabases() (string, error) {
	// MySQL query to list databases as JSON-like string
	// We'll use a simple query and format it if needed,
	// but to keep it consistent with Postgres (json_agg),
	// we'll fetch names and our tool handler expects a JSON array string.
	query := "SELECT JSON_ARRAYAGG(schema_name) FROM information_schema.schemata WHERE schema_name NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys');"
	return m.ExecuteStatement("information_schema", query)
}
