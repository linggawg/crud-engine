package handler

import (
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
	echoGroup.GET("/:table", h.Get)
	echoGroup.POST("/:table", h.Post)
	echoGroup.PUT("/:table/:id", h.Put)
	echoGroup.PATCH("/:table/:id", h.Put)
	echoGroup.DELETE("/:table/:id", h.Delete)
}
