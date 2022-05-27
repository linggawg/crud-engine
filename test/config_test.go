package test

import (
	conn "crud-engine/pkg/database"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	HTTPPort     uint16
	SQLXDatabase conn.SQLXConfig
}

// GlobalEnv global environment
var GlobalEnv Env

func init() {
	errEnv := godotenv.Load("test.env")
	if errEnv != nil {
		panic("Failed to load env file. Make sure .env file is exists!")
	}

	port, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		if os.Getenv("PORT") != "" {
			port, err = strconv.Atoi(os.Getenv("PORT"))
		}
		if err != nil {
			panic("missing HTTP_PORT environment")
		}
	}
	GlobalEnv.HTTPPort = uint16(port)

	GlobalEnv.SQLXDatabase = conn.SQLXConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Dialect:  os.Getenv("DB_DIALECT"),
	}
}
