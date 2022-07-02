package databases

import (
	"engine/bin/config"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

// Init the database from config to database connection
func InitSqlx() *sqlx.DB {
	var (
		toDNS string
	)
	cfg := SQLXConfig{
		Host:     config.GlobalEnv.DBHost,
		Port:     config.GlobalEnv.DBPort,
		Name:     config.GlobalEnv.DBName,
		Username: config.GlobalEnv.DBUser,
		Password: config.GlobalEnv.DBPassword,
		Dialect:  config.GlobalEnv.DBDialect,
		SSLMode:  config.GlobalEnv.DBSSLMode,
	}

	switch cfg.Dialect {
	case "mysql":
		toDNS = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	case "postgres":
		toDNS = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
	default:
		toDNS = ""
	}
	if toDNS == "" {
		panic("Database support only mysql / postgres")
	}
	db, err := sqlx.Connect(cfg.Dialect, toDNS)
	if err != nil {
		panic(fmt.Sprintf("Failed to create pool connection database %s", cfg.Dialect))
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	return db
}

func InitDbs(cfg SQLXConfig) (*sqlx.DB, error) {
	var (
		toDNS string
		err   error
	)
	switch cfg.Dialect {
	case "mysql":
		toDNS = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	case "postgres":
		toDNS = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	default:
		toDNS = ""
	}
	if toDNS == "" {
		return nil, fmt.Errorf("database support only mysql / postgres")
	}
	db, err := sqlx.Connect(cfg.Dialect, toDNS)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	return db, nil
}
