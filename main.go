package main

import (
	"crud-engine/config"
	"crud-engine/handler"
	"crud-engine/mongocontroller"
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

	//  Test Connection Mongo DB
	_, err = config.InitMongo()
	if err != nil {
		panic(err)
	}
	log.Println("MongoDB successfully initialized")

	e := echo.New()

	e.GET("/:table", handler.Get)
	e.POST("/:table", handler.Post)
	e.PUT("/:table/:id", handler.Put)
	e.PATCH("/:table/:id", handler.Put)
	e.DELETE("/:table/:id", handler.Delete)

	e.GET("/getAllUsers", mongocontroller.GetAllUsers)
	e.POST("/createProfile", mongocontroller.CreateProfile)
	e.POST("/getUserProfile", mongocontroller.GetUserProfile)
	e.PUT("/updateProfile/:id", mongocontroller.UpdateProfile)
	e.DELETE("/deleteProfile/:id", mongocontroller.DeleteProfile)

	log.Println("Webserver successfully started")
	log.Println("Listening to port ", os.Getenv("WEBSERVER_LISTEN_ADDRESS"))

	e.Logger.Fatal(e.Start(os.Getenv("WEBSERVER_LISTEN_ADDRESS")))
}
