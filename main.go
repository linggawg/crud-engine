package main

import (
	"crud-engine/config"
	_ "crud-engine/docs"
	"crud-engine/handler"
	"crud-engine/mongocontroller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// @title Echo Swagger Example API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http
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

	// Echo instance
	e := echo.New()

	//// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/", HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

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

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags HealthCheck
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
