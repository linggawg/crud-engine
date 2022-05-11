package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Init() {
	cfg := mysql.Config{
		User:  	os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
}

func CreateCon() *sql.DB {
	return db
}