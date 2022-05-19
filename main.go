package main

import (
	"crud-engine/config"
	_ "crud-engine/docs"
	"crud-engine/handler"
	"crud-engine/mongocontroller"
	conn "crud-engine/pkg/database"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"net/http"
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
	sqlx, err := conn.InitSqlx(config.GlobalEnv.SQLXDatabase)
	if err != nil {
		panic(err)
	}
	log.Println("Database successfully initialized")

	//  Test Connection Mongo DB
	mongoDb, err := conn.InitMongo(config.GlobalEnv.MongoDb)
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

	//initiate user http handler
	crudGroup := e.Group("/sql")
	crudSQLX := handler.New(sqlx)
	crudSQLX.Mount(crudGroup)

	//initiate mongo http handler
	crudMongo := e.Group("/mongodb")
	handlerMongo := mongocontroller.New(mongoDb)
	handlerMongo.Mount(crudMongo)

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
