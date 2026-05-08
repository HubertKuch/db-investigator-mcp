package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type DBDriver interface {
	ExecuteStatement(dbname string, statement string) (string, error)
	ExtractDDL(dbname string) (string, error)
	ListDatabases() (string, error)
}

type PostgresDriver struct{}

func (p *PostgresDriver) ExecuteStatement(dbname string, statement string) (string, error) {
	host, user, password, port, envErr := getDatabaseENV()

	if envErr != nil {
		return "", envErr
	}

	cmd := exec.Command("psql",
		"-h", host,
		"-p", port,
		"-U", user,
		"-d", dbname,
		"-c", statement,
		"-q", "",
		"-t", "",
	)

	cmd.Env = append(os.Environ(), "PGPASSWORD="+password)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	cmdErr := cmd.Run()
	if cmdErr != nil {
		return "", fmt.Errorf("błąd podczas wykonywania zapytania: %s. Szczegóły: %s", cmdErr, stderr.String())
	}

	return out.String(), nil
}

func (p *PostgresDriver) ExtractDDL(dbname string) (string, error) {
	host, user, password, port, envErr := getDatabaseENV()

	if envErr != nil {
		return "", envErr
	}

	cmd := exec.Command("pg_dump", "-h", host, "-p", port, "-U", user, "-s", dbname)

	cmd.Env = append(os.Environ(), "PGPASSWORD="+password)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	cmdErr := cmd.Run()
	if cmdErr != nil {
		return "", fmt.Errorf("błąd podczas pobierania DDL: %s. Szczegóły bazy: %s", cmdErr, stderr.String())
	}

	return out.String(), nil
}

func (p *PostgresDriver) ListDatabases() (string, error) {
	query := "SELECT json_agg(datname) FROM pg_database WHERE datistemplate = false AND datname != 'postgres';"
	return p.ExecuteStatement("postgres", query)
}

func GetDriver() (DBDriver, error) {
	driverType := os.Getenv("DB_TYPE")
	if driverType == "" || driverType == "postgres" {
		return &PostgresDriver{}, nil
	}
	return nil, fmt.Errorf("unsupported database type: %s", driverType)
}

func getDatabaseENV() (string, string, string, string, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")

	if host == "" || user == "" || password == "" || port == "" {
		return "", "", "", "", fmt.Errorf("brak wymaganych zmiennych środowiskowych: DB_HOST, DB_USER, DB_PASSWORD, DB_PORT")
	}
	return host, user, password, port, nil
}
