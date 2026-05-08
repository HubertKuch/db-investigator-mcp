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

	// Test sqlite
	os.Setenv("DB_TYPE", "sqlite")
	os.Setenv("DB_PATH", "/tmp")
	driver, err = GetDriver()
	if err != nil {
		t.Errorf("expected no error for sqlite driver, got %v", err)
	}
	if _, ok := driver.(*SqliteDriver); !ok {
		t.Errorf("expected SqliteDriver for sqlite, got %T", driver)
	}

	// Test mysql
	os.Setenv("DB_TYPE", "mysql")
	driver, err = GetDriver()
	if err != nil {
		t.Errorf("expected no error for mysql driver, got %v", err)
	}
	if _, ok := driver.(*MysqlDriver); !ok {
		t.Errorf("expected MysqlDriver for mysql, got %T", driver)
	}

	// Test unsupported
	os.Setenv("DB_TYPE", "oracle")
	_, err = GetDriver()
	if err == nil {
		t.Error("expected error for unsupported driver, got nil")
	}
}

func TestLoadConfig_Missing(t *testing.T) {
	os.Clearenv()
	_, err := LoadConfig("postgres")
	if err == nil {
		t.Error("expected error for missing ENV, got nil")
	}
}
