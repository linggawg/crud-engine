package main

import (
	"crud-engine/config"
	"crud-engine/handler"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Failed to load env file. Make sure .env file is exists!")
	}
	config.Init()

	e := echo.New()

	e.GET("/"+os.Getenv("TABLE"), handler.Read)
	e.POST("/"+os.Getenv("TABLE"), handler.Create)
	e.DELETE("/"+os.Getenv("TABLE"+"/:id"), handler.Delete)
	e.Logger.Fatal(e.Start(os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT")))
}