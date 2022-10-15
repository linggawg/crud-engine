package usecases

import (
	"context"
	"engine/bin/modules/engine/helpers"
	models "engine/bin/modules/engine/models/domain"
	queriesQuery "engine/bin/modules/queries/repositories/queries"
	servicesQuery "engine/bin/modules/services/repositories/queries"
	"engine/bin/pkg/databases/connection"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/utils"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type EngineQueryUsecase struct {
	queriesQuery  queriesQuery.QueriesPostgre
	servicesQuery servicesQuery.ServicesPostgre
	db            connection.Connection
	repo          *helpers.BulkRepository
}

func NewQueryUsecase(queriesQuery queriesQuery.QueriesPostgre, servicesQuery servicesQuery.ServicesPostgre, db connection.Connection, repo helpers.BulkRepository) *EngineQueryUsecase {
	return &EngineQueryUsecase{
		queriesQuery:  queriesQuery,
		servicesQuery: servicesQuery,
		db:            db,
		repo:          &repo,
	}
}

func (h *EngineQueryUsecase) Get(ctx context.Context, engineConfig models.EngineConfig, table string, payload *models.GetList) (result utils.Result) {
	var (
		sqlStatement string
		sqlCount     string
		repository   = h.repo.GetBulkRepository(engineConfig.Dbs.Dialect)
	)

	tx, err := h.db.GetBulkConnectionSql(engineConfig.Dbs)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	if payload.Key != "" {
		queriesModels, err := h.queriesQuery.FindOneByKey(ctx, payload.Key)
		if err != nil {
			errObj := httpError.Conflict{}
			errObj.Message = "key did not exist"
			result.Error = errObj
			return result
		}

		filter := payload.Filter
		if filter != "" {
			filter = " WHERE " + filter
		}

		sqlStatement = queriesModels.QueryDefinition + filter + " "
		sqlCount = sqlStatement

		servicesModels, _ := h.servicesQuery.FindOneByQueryID(ctx, queriesModels.ID)
		primaryKey, _ := repository.FindPrimaryKey(ctx, tx, *servicesModels.ServiceUrl)
		if primaryKey != nil {
			sqlStatement = setQueryPagination(primaryKey, sqlStatement, payload)
		}
	} else {
		primaryKey, err := repository.FindPrimaryKey(ctx, tx, table)
		if err != nil {
			errObj := httpError.NewNotFound()
			errObj.Message = "Undefined primary key"
			result.Error = errObj
			return result
		}
		isDistinct := ""
		if payload.IsDistinct {
			isDistinct = "DISTINCT "
		}

		filter := payload.Filter
		if filter != "" {
			filter = "WHERE " + filter
		}

		columns := payload.Columns
		if columns == "" {
			columns = "*"
		}

		sqlStatement = fmt.Sprintf(utils.QueryGet, isDistinct, columns, table, filter)
		sqlCount = sqlStatement
		sqlStatement = setQueryPagination(primaryKey, sqlStatement, payload)
	}

	totalItems, err := repository.CountData(ctx, tx, sqlCount)
	if err != nil {
		errObj := httpError.NewUnprocessableEntity()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	tableData, err := repository.FindData(ctx, tx, sqlStatement)
	if err != nil {
		errObj := httpError.NewUnprocessableEntity()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	if len(engineConfig.ResourcesMappingList) > 0 {
		for i, data := range tableData {
			tableData[i] = helpers.ConvertToSourceAlias(data, engineConfig.ResourcesMappingList)
		}
	}

	result.Data = tableData
	result.MetaData = map[string]interface{}{
		"page": func() *int {
			if payload.Page != nil {
				if *payload.Page == 0 {
					*payload.Page = 1
				}
			}
			return payload.Page
		}(),
		"quantity": len(tableData),
		"totalPage": func() *float64 {
			if payload.Size != nil {
				maxPage := math.Ceil(float64(totalItems) / float64(*payload.Size))
				return &maxPage
			} else {
				return nil
			}
		}(),
		"totalData": totalItems,
	}
	return result
}

func setQueryPagination(pk *models.PrimaryKey, query string, p *models.GetList) (newQuery string) {
	var (
		pagination string
		page       int
		size       int
	)
	if p.Sort != "" {
		query += fmt.Sprintf(" ORDER BY %s ", p.Sort)
	}
	if p.Size != nil {
		size = *p.Size
		if p.Page != nil {
			page = (*p.Page - 1) * size
			if page < 0 {
				page = 0
			}
		}

		if p.Sort == "" && strings.EqualFold(pk.Format, "int") {
			if p.Filter != "" {
				pagination = fmt.Sprintf("AND "+utils.QueryPaginationOrdinal, pk.Column, strconv.Itoa(page), pk.Column, strconv.Itoa(size))
			} else {
				pagination = fmt.Sprintf("WHERE "+utils.QueryPaginationOrdinal, pk.Column, strconv.Itoa(page), pk.Column, strconv.Itoa(size))
			}
		} else {
			pagination = fmt.Sprintf(utils.QueryPaginationDefault, strconv.Itoa(size), strconv.Itoa(page))
		}

		return query + pagination
	}
	return query
}
