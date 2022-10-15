package usecases

import (
	"context"
	"database/sql"
	"engine/bin/modules/dbs/repositories/queries"
	engineModels "engine/bin/modules/engine/models/domain"
	queriesQuery "engine/bin/modules/queries/repositories/queries"
	resourcesMapping "engine/bin/modules/resources-mapping/repositories/queries"
	serviceModels "engine/bin/modules/services/models/domain"
	servicesQuery "engine/bin/modules/services/repositories/queries"
	usersServicesQuery "engine/bin/modules/users-services/repositories/queries"
	usersQuery "engine/bin/modules/users/repositories/queries"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/utils"
	"fmt"
	"strings"
)

type DbsQueryUsecase struct {
	dbsPostgreQuery    queries.DbsPostgre
	servicesQuery      servicesQuery.ServicesPostgre
	usersServicesQuery usersServicesQuery.UsersServicesPostgre
	usersQuery         usersQuery.UsersPostgre
	queriesQuery       queriesQuery.QueriesPostgre
	resourcesMapping   resourcesMapping.ResourcesMappingPostgre
}

func NewQueryUsecase(dbsPostgreQuery queries.DbsPostgre, queriesQuery queriesQuery.QueriesPostgre, servicesQuery servicesQuery.ServicesPostgre, usersServicesQuery usersServicesQuery.UsersServicesPostgre, usersQuery usersQuery.UsersPostgre, resourcesMapping resourcesMapping.ResourcesMappingPostgre) *DbsQueryUsecase {
	return &DbsQueryUsecase{
		dbsPostgreQuery:    dbsPostgreQuery,
		servicesQuery:      servicesQuery,
		usersServicesQuery: usersServicesQuery,
		usersQuery:         usersQuery,
		queriesQuery:       queriesQuery,
		resourcesMapping:   resourcesMapping,
	}
}

func (h *DbsQueryUsecase) GetDbsConnection(ctx context.Context, userId, method, serviceUrl, key string) utils.Result {
	var (
		result  utils.Result
		service *serviceModels.Services
		err     error
	)

	//Services validation
	if strings.EqualFold(serviceUrl, "query") {
		//Queries validation
		queries, err := h.queriesQuery.FindOneByKey(ctx, key)
		if err != nil {
			if err == sql.ErrNoRows {
				errObj := httpError.NewNotFound()
				errObj.Message = "data query definition tidak ditemukan"
				result.Error = errObj
			} else {
				errObj := httpError.NewInternalServerError()
				errObj.Message = err.Error()
				result.Error = errObj
			}
			return result
		}

		service, err = h.servicesQuery.FindOneByQueryID(ctx, queries.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				errObj := httpError.NewNotFound()
				errObj.Message = fmt.Sprintf("data service dengan key: %s tidak ditemukan", key)
				result.Error = errObj
			} else {
				errObj := httpError.NewInternalServerError()
				errObj.Message = err.Error()
				result.Error = errObj
			}
			return result
		}
	} else {
		service, err = h.servicesQuery.FindOneByServiceUrlAndMethodAndQueryIsNull(ctx, serviceUrl, method)
		if err != nil {
			if err == sql.ErrNoRows {
				errObj := httpError.NewNotFound()
				errObj.Message = "data service url tidak ditemukan"
				result.Error = errObj
			} else {
				errObj := httpError.NewInternalServerError()
				errObj.Message = err.Error()
				result.Error = errObj
			}
			return result
		}
	}

	//Users Services validation
	_, err = h.usersServicesQuery.FindOneByServiceIDAndUserId(ctx, service.ID, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			errObj := httpError.NewNotFound()
			errObj.Message = "data users services tidak ditemukan"
			result.Error = errObj
		} else {
			errObj := httpError.NewInternalServerError()
			errObj.Message = err.Error()
			result.Error = errObj
		}
		return result
	}

	//Resources Mapping validation
	resourcesMapping, err := h.resourcesMapping.FindByServiceId(ctx, service.ID)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	//Dbs validation
	dbs, err := h.dbsPostgreQuery.FindOneByID(ctx, service.DbID)
	if err != nil {
		if err == sql.ErrNoRows {
			errObj := httpError.NewNotFound()
			errObj.Message = "data dbs tidak ditemukan"
			result.Error = errObj
		} else {
			errObj := httpError.NewInternalServerError()
			errObj.Message = err.Error()
			result.Error = errObj
		}
		return result
	}

	result.Data = engineModels.EngineConfig{
		Dbs:                  *dbs,
		ResourcesMappingList: resourcesMapping,
	}
	return result
}
