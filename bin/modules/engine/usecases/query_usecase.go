package usecases

import (
	"context"
	dbsModels "engine/bin/modules/dbs/models/domain"
	"engine/bin/modules/engine/helpers"
	models "engine/bin/modules/engine/models/domain"
	"engine/bin/modules/engine/repositories/queries"
	"engine/bin/pkg/databases"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/utils"
	"math"
	"net/url"
	"strconv"
)

type engineQueryUsecase struct {
	engineQuery queries.EngineSQL
	db          *databases.BulkConnectionPkg
}

func NewQueryUsecase(engineQuery queries.EngineSQL, db databases.BulkConnectionPkg) *engineQueryUsecase {
	return &engineQueryUsecase{
		engineQuery: engineQuery,
		db:          &db,
	}
}

func (h *engineQueryUsecase) Get(ctx context.Context, dbs dbsModels.Dbs, table string, payload *models.GetList) (result utils.Result) {
	var (
		sqlStatement string
		sqlCount     string
	)

	tx, err := h.db.GetBulkConnectionSql(dbs.ID)
	if err != nil {
		tx, err = helpers.CreateConnection(dbs)
		if err != nil {
			errObj := httpError.InternalServerError{}
			errObj.Message = err.Error()
			result.Error = errObj
			return result
		}
		err = h.db.AddBulkConnectionSql(dbs.ID, tx)
		if err != nil {
			errObj := httpError.Conflict{}
			errObj.Message = "cannot add connection to bulk manager connection"
			result.Error = errObj
			return result
		}
	}

	if payload.IsQuery {
		table, _ = url.QueryUnescape(table)
		sqlStatement = table
		sqlCount = table
	} else {
		primaryKey, err := h.engineQuery.FindPrimaryKey(ctx, tx, dbs.Dialect, table)
		if err != nil {
			errObj := httpError.NewNotFound()
			errObj.Message = "Primary key tidak ditemukan"
			result.Error = errObj
			return result
		}
		isDistinct := ""
		if payload.IsDistinct {
			isDistinct = "DISTINCT "
		}

		query := payload.Query
		if query != "" {
			query = " WHERE " + query
		}

		colls := payload.Colls
		if colls == "" {
			colls = "*"
		}

		sqlStatement = "SELECT " + isDistinct + colls + " FROM " + table + query
		sqlCount = sqlStatement
		sqlStatement = setQueryPagination(sqlStatement, primaryKey.Column, payload)
	}

	totalItems, err := h.engineQuery.CountData(ctx, tx, sqlCount)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		result.Error = errObj
		return result
	}

	tableData, err := h.engineQuery.FindData(ctx, tx, sqlStatement)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		result.Error = errObj
		return result
	}

	result.Data = tableData
	result.MetaData = map[string]interface{}{
		"page":     payload.Page,
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

func setQueryPagination(query string, primaryKey string, p *models.GetList) (newQuery string) {
	if p != nil {
		if p.Sort != "" {
			query += " ORDER BY " + p.Sort
		} else {
			query += " ORDER BY " + primaryKey + " ASC"
		}
		if p.Size != nil {
			sz := *p.Size
			size := strconv.Itoa(sz)
			query += " LIMIT " + size
			if p.Page != nil {
				pg := *p.Page
				page := (pg - 1) * sz
				if page > 0 {
					query += " OFFSET " + strconv.Itoa(page)
				}
			}
		}
	}
	return query
}
