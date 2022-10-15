package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Env struct {
	RootApp          string
	HTTPPort         uint16
	EngineApiUrl     string
	StrictRestfulAPI bool
	APISecret        string
	DBHost           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBPort           uint16
	DBSSLMode        string
	DBDialect        string
	EngineUser       string
	EnginePassword   string
	EngineRole       string
}

// GlobalEnv global environment
var GlobalEnv Env

const ProjectDirName = "agree-logtan-engine" // change to relevant project name

func loadEnv() {
	projectName := regexp.MustCompile(`^(.*` + ProjectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Println(err)
	}
	rootApp := strings.TrimSuffix(currentWorkDirectory, "/bin/config")
	os.Setenv("APP_PATH", rootApp)
	GlobalEnv.RootApp = rootApp
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

	GlobalEnv.EngineApiUrl, ok = os.LookupEnv("ENGINE_API_URL")
	if !ok {
		panic("missing ENGINE_API_URL environment")
	}

	strictRestful, ok := os.LookupEnv("STRICT_RESTFUL_API")
	if !ok {
		panic("missing STRICT_RESTFUL_API environment")
	}
	GlobalEnv.StrictRestfulAPI, err = strconv.ParseBool(strictRestful)
	if err != nil {
		panic("failed parsing StrictRestfulAPI")
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

	GlobalEnv.EngineUser, ok = os.LookupEnv("ENGINE_USER_ADMIN")
	if !ok {
		panic("missing ENGINE_USER_ADMIN environment")
	}

	GlobalEnv.EnginePassword, ok = os.LookupEnv("ENGINE_PASSWORD_ADMIN")
	if !ok {
		panic("missing ENGINE_PASSWORD_ADMIN environment")
	}

	GlobalEnv.EngineRole, ok = os.LookupEnv("ENGINE_ROLE")
	if !ok {
		panic("missing ENGINE_ROLEs environment")
	}

}
