package utils

import (
	"os"
	"testing"
)

func TestGetDriver(t *testing.T) {
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

func TestGetDatabaseENV_Missing(t *testing.T) {
	os.Clearenv()
	_, _, _, _, err := getDatabaseENV()
	if err == nil {
		t.Error("expected error for missing ENV, got nil")
	}
}
