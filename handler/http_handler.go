package handler

import (
	"crud-engine/pkg/middleware"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type HttpSqlx struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *HttpSqlx {
	return &HttpSqlx{db}
}

// Mount function
func (h *HttpSqlx) Mount(echoGroup *echo.Group) {
	//for development
	echoGroup.GET("/dev/:table", h.Get)
	echoGroup.POST("/dev/:table", h.Post)
	echoGroup.PUT("/dev/:table/:value", h.Put)
	echoGroup.PATCH("/dev/:table/:value", h.Put)
	echoGroup.DELETE("/dev/:table/:value", h.Delete)

	echoGroup.GET("/:table", h.Get, middleware.VerifyBearer())
	echoGroup.POST("/:table", h.Post, middleware.VerifyBearer())
	echoGroup.PUT("/:table/:value", h.Put, middleware.VerifyBearer())
	echoGroup.PATCH("/:table/:value", h.Put, middleware.VerifyBearer())
	echoGroup.DELETE("/:table/:value", h.Delete, middleware.VerifyBearer())
}
