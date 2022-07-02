package main

import (
	"engine/bin/config"
	_ "engine/bin/docs"
	"engine/bin/modules/engine/handlers"
	usersHandler "engine/bin/modules/users/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func init() {
	log.SetPrefix("[API-CRUD ENGINE SERVICE] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

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

// @securityDefinitions.apikey Authorization
// @in header
// @name Authorization
func main() {
	// Echo instance
	e := echo.New()

	//// Middleware
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/", HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	//initiate auth http handler
	usersHandler.New().Mount(e.Group(""))

	//initiate user http handler
	handlers.New().Mount(e.Group("/engine"))

	listenerPort := fmt.Sprintf(":%d", config.GlobalEnv.HTTPPort)
	log.Println("Webserver successfully started")
	log.Println("Listening to port ", config.GlobalEnv.HTTPPort)
	e.Logger.Fatal(e.Start(listenerPort))
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
