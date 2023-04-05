package main

import (
	"engine/bin/config"
	engineHandler "engine/bin/modules/engine/handlers"
	servicesHandler "engine/bin/modules/services/handlers"
	usersServicesHandler "engine/bin/modules/users-services/handlers"
	usersHandler "engine/bin/modules/users/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	log.SetPrefix("[API-CRUD ENGINE SERVICE] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/engine/", HealthCheck)

	engineGroup := e.Group("/engine")

	servicesHTTP := servicesHandler.New()
	servicesHTTP.Mount(engineGroup)

	usersHTTP := usersHandler.New()
	usersHTTP.Mount(engineGroup)

	usersServicesHTTP := usersServicesHandler.New()
	usersServicesHTTP.Mount(engineGroup)

	engineHTTP := engineHandler.New()
	engineHTTP.Mount(engineGroup)

	listenerPort := fmt.Sprintf(":%d", config.GlobalEnv.HTTPPort)
	e.Logger.Fatal(e.Start(listenerPort))
}

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
