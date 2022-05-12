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
	_, err := config.Init()
	if err != nil {
		panic(err)
	}
	log.Println("Database successfully initialized")

	e := echo.New()

	e.GET("/:table", handler.Read)
	e.POST("/:table", handler.Create)
	e.DELETE("/:table/:id", handler.Delete)
	e.Logger.Fatal(e.Start(os.Getenv("WEBSERVER_LISTEN_ADDRESS")))

	log.Println("Webserver successfully started")
	log.Println("Listening to port ", os.Getenv("WEBSERVER_LISTEN_ADDRESS"))

	e.Logger.Fatal(e.Start(os.Getenv("WEBSERVER_LISTEN_ADDRESS")))
}
