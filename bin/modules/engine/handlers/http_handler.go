package handlers

import (
	"context"
	"engine/bin/middleware"
	dbsQueries "engine/bin/modules/dbs/repositories/queries"
	dbsUsecase "engine/bin/modules/dbs/usecases"
	"engine/bin/modules/engine/helpers"
	engineUsecase "engine/bin/modules/engine/usecases"
	queriesQueries "engine/bin/modules/queries/repositories/queries"
	resourcesMapping "engine/bin/modules/resources-mapping/repositories/queries"
	servicesQueries "engine/bin/modules/services/repositories/queries"
	usersServicesQueries "engine/bin/modules/users-services/repositories/queries"
	usersQueries "engine/bin/modules/users/repositories/queries"
	"engine/bin/pkg/databases"
	"engine/bin/pkg/databases/connection"
	"engine/bin/pkg/utils"

	"github.com/labstack/echo/v4"
)

type EngineHTTPHandler struct {
	EngineQueryUsecase   engineUsecase.QueryUsecase
	EngineCommandUsecase engineUsecase.CommandUsecase
	DbsQueryUsecase      dbsUsecase.QueryUsecase
}

func New() *EngineHTTPHandler {
	db := databases.InitSqlx()
	bulkConnection := connection.InitBulkConnectionPkg()
	bulkRepository := helpers.InitBulkRepository()

	servicesPostgreQuery := servicesQueries.NewServicesQuery(db)
	usersServicesPostgreQuery := usersServicesQueries.NewUsersServicesQuery(db)
	usersQuery := usersQueries.NewUsersQuery(db)
	queriesQuery := queriesQueries.NewQueriesQuery(db)
	resourcesMaping := resourcesMapping.NewResourcesMappingQuery(db)

	engineQueryUsecase := engineUsecase.NewQueryUsecase(queriesQuery, servicesPostgreQuery, bulkConnection, *bulkRepository)
	engineCommandUsecase := engineUsecase.NewCommandUsecase(bulkConnection, *bulkRepository)

	dbsPostgreQuery := dbsQueries.NewDbsQuery(db)
	dbsQueryUsecase := dbsUsecase.NewQueryUsecase(dbsPostgreQuery, queriesQuery, servicesPostgreQuery, usersServicesPostgreQuery, usersQuery, resourcesMaping)

	engineCommandUsecase.SetupDataset(context.Background(), utils.DialectPostgres, db)

	return &EngineHTTPHandler{
		DbsQueryUsecase:      dbsQueryUsecase,
		EngineQueryUsecase:   engineQueryUsecase,
		EngineCommandUsecase: engineCommandUsecase,
	}
}

func (h *EngineHTTPHandler) Mount(echoGroup *echo.Group) {
	echoGroup.GET("/v1/:table", h.Get, middleware.VerifyBearer())
	echoGroup.POST("/v1/:table", h.Post, middleware.VerifyBearer())
	echoGroup.PUT("/v1/:table/:value", h.Put, middleware.VerifyBearer())
	echoGroup.PATCH("/v1/:table/:value", h.Patch, middleware.VerifyBearer())
	echoGroup.DELETE("/v1/:table/:value", h.Delete, middleware.VerifyBearer())
}
