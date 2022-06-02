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
	echoGroup.GET("/:table", h.Get, middleware.VerifyBearer())
	echoGroup.POST("/:table", h.Post, middleware.VerifyBearer())
	echoGroup.PUT("/:table/:value", h.Put, middleware.VerifyBearer())
	echoGroup.PATCH("/:table/:value", h.Put, middleware.VerifyBearer())
	echoGroup.DELETE("/:table/:value", h.Delete, middleware.VerifyBearer())
}
