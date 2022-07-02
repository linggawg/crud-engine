package handlers

import (
	"engine/bin/middleware"
	dbsQueries "engine/bin/modules/dbs/repositories/queries"
	dbsUsecase "engine/bin/modules/dbs/usecases"
	engineCommands "engine/bin/modules/engine/repositories/commands"
	engineQueries "engine/bin/modules/engine/repositories/queries"
	engineUsecase "engine/bin/modules/engine/usecases"
	servicesQueries "engine/bin/modules/services/repositories/queries"
	usersQueries "engine/bin/modules/users/repositories/queries"
	userServicesQueries "engine/bin/modules/userservice/repositories/queries"
	"engine/bin/pkg/databases"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type EngineHTTPHandler struct {
	db                   *sqlx.DB
	engineQueryUsecase   engineUsecase.QueryUsecase
	engineCommandUsecase engineUsecase.CommandUsecase
	dbsQueryUsecase      dbsUsecase.QueryUsecase
}

func New() *EngineHTTPHandler {
	db := databases.InitSqlx()
	bulkConnection := databases.InitBulkConnectionPkg()

	engineQuery := engineQueries.NewEngineQuery()
	engineCommand := engineCommands.NewEngineCommand()
	engineQueryUsecase := engineUsecase.NewQueryUsecase(engineQuery, *bulkConnection)
	engineCommandUsecase := engineUsecase.NewCommandUsecase(engineCommand, engineQuery, *bulkConnection)

	servicesPostgreQuery := servicesQueries.NewServicesQuery(db)
	userServicePostgreQuery := userServicesQueries.NewUserServiceQuery(db)
	usersQuery := usersQueries.NewUsersQuery(db)

	dbsPostgreQuery := dbsQueries.NewDbsQuery(db)
	dbsQueryUsecase := dbsUsecase.NewQueryUsecase(dbsPostgreQuery, servicesPostgreQuery, userServicePostgreQuery, usersQuery)

	return &EngineHTTPHandler{
		db:                   db,
		dbsQueryUsecase:      dbsQueryUsecase,
		engineQueryUsecase:   engineQueryUsecase,
		engineCommandUsecase: engineCommandUsecase,
	}
}

// Mount function
func (h *EngineHTTPHandler) Mount(echoGroup *echo.Group) {
	echoGroup.GET("/:table", h.Get, middleware.VerifyBearer())
	echoGroup.POST("/:table", h.Post, middleware.VerifyBearer())
	echoGroup.PUT("/:table/:value", h.Put, middleware.VerifyBearer())
	echoGroup.DELETE("/:table/:value", h.Delete, middleware.VerifyBearer())
}
