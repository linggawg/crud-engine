package main

import (
	"engine/bin/config"
	_ "engine/bin/docs"
	engineHandler "engine/bin/modules/engine/handlers"
	servicesHandler "engine/bin/modules/services/handlers"
	usersServicesHandler "engine/bin/modules/users-services/handlers"
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

// @title Echo Swagger Engine Services
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /engine/

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
	e.GET("/engine/", HealthCheck)
	e.GET("/engine/swagger/*", echoSwagger.WrapHandler)

	//initiate auth http handler
	engineGroup := e.Group("/engine")

	//initiate services http handler
	servicesHTTP := servicesHandler.New()
	servicesHTTP.Mount(engineGroup)

	//initiate users http handler
	usersHTTP := usersHandler.New()
	usersHTTP.Mount(engineGroup)

	//initiate services http handler
	usersServicesHTTP := usersServicesHandler.New()
	usersServicesHTTP.Mount(engineGroup)

	//initiate engine http handler
	engineHTTP := engineHandler.New()
	engineHTTP.Mount(engineGroup)

	listenerPort := fmt.Sprintf(":%d", config.GlobalEnv.HTTPPort)
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
