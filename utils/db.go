package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func ExecuteStatement(dbname string, statement string) (string, error) {
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

func ExtractDDL(dbname string) (string, error) {
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
