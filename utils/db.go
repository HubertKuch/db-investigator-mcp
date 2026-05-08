package utils

import (
	"fmt"
	"os"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBPath   string // For SQLite
}

type DBDriver interface {
	ExecuteStatement(dbname string, statement string) (string, error)
	ExtractDDL(dbname string) (string, error)
	ListDatabases() (string, error)
}

func GetDriver() (DBDriver, error) {
	driverType := os.Getenv("DB_TYPE")
	if driverType == "" {
		driverType = "postgres"
	}

	cfg, err := LoadConfig(driverType)
	if err != nil {
		return nil, err
	}

	switch driverType {
	case "postgres":
		return &PostgresDriver{cfg: cfg}, nil
	case "sqlite":
		return &SqliteDriver{cfg: cfg}, nil
	case "mysql":
		return &MysqlDriver{cfg: cfg}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", driverType)
	}
}

func LoadConfig(driverType string) (*Config, error) {
	cfg := &Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Port:     os.Getenv("DB_PORT"),
		DBPath:   os.Getenv("DB_PATH"),
	}

	if driverType == "postgres" || driverType == "mysql" {
		if cfg.Host == "" || cfg.User == "" || cfg.Password == "" || cfg.Port == "" {
			return nil, fmt.Errorf("missing required environment variables for %s: DB_HOST, DB_USER, DB_PASSWORD, DB_PORT", driverType)
		}
	}

	return cfg, nil
}
