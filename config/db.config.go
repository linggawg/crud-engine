package config

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

var db *sqlx.DB

//func Init() {
//	cfg := mysql.Config{
//		User:   os.Getenv("DB_USER"),
//		Passwd: os.Getenv("DB_PASSWORD"),
//		Net:    "tcp",
//		Addr:   os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
//		DBName: os.Getenv("DB_NAME"),
//	}
//
//	var err error
//	db, err = sql.Open("mysql", cfg.FormatDSN())
//	if err != nil {
//		log.Fatal(err)
//	}
//}

// Init the database from env to database connection
func Init() (*sqlx.DB, error) {
	// toDNS := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%2FJakarta&charset=utf8&autocommit=false", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	var toDNS string
	switch os.Getenv("DB_DIALECT") {
	case "mysql":
		toDNS = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	case "postgres":
		toDNS = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	default:
		toDNS = ""
	}
	if toDNS == "" {
		return nil, fmt.Errorf("Database support only mysql / postgres")
	}
	var err error
	db, err = sqlx.Connect(os.Getenv("DB_DIALECT"), toDNS)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateCon() *sqlx.DB {
	return db
}
