package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Env struct {
	HTTPPort           uint16
	AccessTokenExpired time.Duration
	APISecret          string
	DBHost             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBPort             uint16
	DBSSLMode          string
	DBDialect          string
}

// GlobalEnv global environment
var GlobalEnv Env

const projectDirName = "crud-engine-be" // change to relevant project name

func loadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Println(err)
	}
}

func init() {
	loadEnv()
	var ok bool

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

	GlobalEnv.AccessTokenExpired, err = time.ParseDuration("100m")
	if err != nil {
		panic("failed parsing AccessTokenExpired")
	}

	GlobalEnv.APISecret, ok = os.LookupEnv("API_SECRET")
	if !ok {
		panic("missing API_SECRET environment")
	}

	GlobalEnv.DBHost, ok = os.LookupEnv("DB_HOST")
	if !ok {
		panic("missing DB_HOST environment")
	}

	if dbPort, err := strconv.Atoi(os.Getenv("DB_PORT")); err != nil {
		panic("missing DB_PORT environment")
	} else {
		GlobalEnv.DBPort = uint16(dbPort)
	}

	GlobalEnv.DBUser, ok = os.LookupEnv("DB_USER")
	if !ok {
		panic("missing DB_USER environment")
	}

	GlobalEnv.DBPassword, ok = os.LookupEnv("DB_PASSWORD")
	if !ok {
		panic("missing DB_PASSWORD environment")
	}

	GlobalEnv.DBName, ok = os.LookupEnv("DB_NAME")
	if !ok {
		panic("missing DB_NAME environment")
	}

	GlobalEnv.DBSSLMode, ok = os.LookupEnv("DB_SSLMODE")
	if !ok {
		panic("missing DB_SSLMODE environment")
	}

	GlobalEnv.DBDialect, ok = os.LookupEnv("DB_DIALECT")
	if !ok {
		panic("missing DB_DIALECT environment")
	}

}
