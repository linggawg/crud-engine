package handler

import (
	dbsQueries "crud-engine/modules/dbs/repositories/queries"
	dbsUsecase "crud-engine/modules/dbs/usecases"
	"crud-engine/pkg/middleware"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type HttpSqlx struct {
	db              *sqlx.DB
	dbsQueryUsecase dbsUsecase.QueryUsecase
}

func New(db *sqlx.DB) *HttpSqlx {

	dbsPostgreQuery := dbsQueries.NewDbsQuery(db)
	dbsQueryUsecase := dbsUsecase.NewQueryUsecase(dbsPostgreQuery)

	return &HttpSqlx{
		db:              db,
		dbsQueryUsecase: dbsQueryUsecase,
	}
}

// Mount function
func (h *HttpSqlx) Mount(echoGroup *echo.Group) {
	echoGroup.GET("/:table", h.Get, middleware.VerifyBearer())
	echoGroup.POST("/:table", h.Post, middleware.VerifyBearer())
	echoGroup.PUT("/:table/:value", h.Put, middleware.VerifyBearer())
	echoGroup.PATCH("/:table/:value", h.Put, middleware.VerifyBearer())
	echoGroup.DELETE("/:table/:value", h.Delete, middleware.VerifyBearer())
}
