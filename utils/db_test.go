package utils

import (
	"os"
	"testing"
)

func TestGetDriver(t *testing.T) {
	// Setup env for success
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "pass")
	os.Setenv("DB_PORT", "5432")

	// Test default (postgres)
	os.Setenv("DB_TYPE", "")
	driver, err := GetDriver()
	if err != nil {
		t.Errorf("expected no error for default driver, got %v", err)
	}
	if _, ok := driver.(*PostgresDriver); !ok {
		t.Errorf("expected PostgresDriver for default, got %T", driver)
	}

	// Test explicit postgres
	os.Setenv("DB_TYPE", "postgres")
	driver, err = GetDriver()
	if err != nil {
		t.Errorf("expected no error for postgres driver, got %v", err)
	}
	if _, ok := driver.(*PostgresDriver); !ok {
		t.Errorf("expected PostgresDriver for explicit postgres, got %T", driver)
	}

	// Test unsupported
	os.Setenv("DB_TYPE", "mysql")
	_, err = GetDriver()
	if err == nil {
		t.Error("expected error for unsupported driver, got nil")
	}
}

func TestLoadConfig_Missing(t *testing.T) {
	os.Clearenv()
	_, err := LoadConfig()
	if err == nil {
		t.Error("expected error for missing ENV, got nil")
	}
}
