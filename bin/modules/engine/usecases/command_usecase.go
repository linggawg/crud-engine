package usecases

import (
	"context"
	dbsModels "engine/bin/modules/dbs/models/domain"
	"engine/bin/modules/engine/helpers"
	models "engine/bin/modules/engine/models/domain"
	"engine/bin/modules/engine/repositories/commands"
	"engine/bin/modules/engine/repositories/queries"
	"engine/bin/pkg/databases"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/utils"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type engineCommandUsecase struct {
	engineCommand commands.EngineSQL
	engineQuery   queries.EngineSQL
	db            *databases.BulkConnectionPkg
}

func NewCommandUsecase(engineCommand commands.EngineSQL, engineQuery queries.EngineSQL, db databases.BulkConnectionPkg) *engineCommandUsecase {
	return &engineCommandUsecase{
		engineCommand: engineCommand,
		engineQuery:   engineQuery,
		db:            &db,
	}
}

func (h *engineCommandUsecase) Insert(ctx context.Context, dbs dbsModels.Dbs, payload *models.EngineRequest) (result utils.Result) {
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

	primaryKey, err := h.engineQuery.FindPrimaryKey(ctx, tx, dbs.Dialect, payload.Table)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Primary key tidak ditemukan"
		result.Error = errObj
		return result
	}
	schema, err := h.engineQuery.SelectInformationSchema(ctx, tx, dbs.Dialect, payload.Table)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = "Internal Server Error"
		result.Error = errObj
		return result
	}
	sql, args, err := h.SqlStatementInsert(dbs.Dialect, primaryKey, payload, schema)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	err = h.engineCommand.InsertOne(ctx, tx, sql, args)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	result.Data = payload.Data
	return result
}

func (h *engineCommandUsecase) SqlStatementInsert(dialect string, primaryKey *models.PrimaryKey, payload *models.EngineRequest, informationSchemas []models.InformationSchema) (sql string, args []interface{}, err error) {
	var columns, values []string
	for key := range payload.Data {
		if key != primaryKey.Column {
			columns = append(columns, key)
			values = append(values, "?")
			if payload.Data[key] == nil {
				for _, i := range informationSchemas {
					if i.IsNullable == "NO" && i.ColumName == key {
						errorMessage := fmt.Sprintf("validation for '%s' failed on the 'required' tag", i.ColumName)
						return "", args, errors.New(errorMessage)
					}
				}
			}
			args = append(args, payload.Data[key])
		}
	}
	if !strings.EqualFold(primaryKey.Format, "int") {
		columns = append(columns, primaryKey.Column)
		values = append(values, "?")
		args = append(args, uuid.New().String())
	}

	for _, i := range informationSchemas {
		if i.IsNullable == "NO" && i.ColumName != primaryKey.Column {
			if !strings.Contains(strings.Join(columns, ","), i.ColumName) {
				errorMessage := fmt.Sprintf(": Error:validation for '%s' failed on the 'required' tag", i.ColumName)
				return "", args, errors.New(errorMessage)
			}
		}
	}
	sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", payload.Table, strings.Join(columns, ","), strings.Join(values, ","))
	return helpers.SetQuery(dialect, sql), args, nil
}

func (h *engineCommandUsecase) Update(ctx context.Context, dbs dbsModels.Dbs, payload *models.EngineRequest) (result utils.Result) {
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

	schema, err := h.engineQuery.SelectInformationSchema(ctx, tx, dbs.Dialect, payload.Table)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = "Internal Server Error"
		result.Error = errObj
		return result
	}
	isFoundField := false
	for _, i := range schema {
		if i.ColumName == payload.FieldId {
			isFoundField = true
			break
		}
	}
	if !isFoundField {
		errObj := httpError.NewInternalServerError()
		errObj.Message = "field_id '" + payload.FieldId + "' is not found"
		result.Error = errObj
		return result
	}

	primaryKey, err := h.engineQuery.FindPrimaryKey(ctx, tx, dbs.Dialect, payload.Table)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Primary key tidak ditemukan"
		result.Error = errObj
		return result
	}
	sql, args, err := h.SqlStatementUpdate(dbs.Dialect, payload, primaryKey, schema)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	err = h.engineCommand.UpdateOne(ctx, tx, sql, args)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	result.Data = payload.Data
	return result
}

func (h *engineCommandUsecase) SqlStatementUpdate(dialect string, payload *models.EngineRequest, primaryKey *models.PrimaryKey, informationSchemas []models.InformationSchema) (sql string, args []interface{}, err error) {
	var setData []string
	for key := range payload.Data {
		if key != primaryKey.Column {
			if payload.Data[key] == nil {
				for _, i := range informationSchemas {
					if i.ColumName == key && i.IsNullable == "NO" {
						err = fmt.Errorf("validation for '%s' failed on the 'required' tag", i.ColumName)
						return "", args, err
					}
				}
			}
			setData = append(setData, fmt.Sprintf("%s=?", key))
			args = append(args, payload.Data[key])
		}
	}
	args = append(args, payload.Value)
	sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", payload.Table, strings.Join(setData, ","), payload.FieldId)
	return helpers.SetQuery(dialect, sql), args, nil
}

func (h *engineCommandUsecase) Delete(ctx context.Context, dbs dbsModels.Dbs, payload *models.EngineRequest) (result utils.Result) {
	var args []interface{}

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

	schema, err := h.engineQuery.SelectInformationSchema(ctx, tx, dbs.Dialect, payload.Table)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = "Internal Server Error"
		result.Error = errObj
		return result
	}
	isFoundField := false
	for _, i := range schema {
		if i.ColumName == payload.FieldId {
			isFoundField = true
			break
		}
	}
	if !isFoundField {
		errObj := httpError.NewInternalServerError()
		errObj.Message = "field_id '" + payload.FieldId + "' is not found"
		result.Error = errObj
		return result
	}
	sql := "DELETE FROM " + payload.Table + " WHERE " + payload.FieldId + " = ?"
	args = append(args, payload.Value)
	err = h.engineCommand.DeleteOne(ctx, tx, helpers.SetQuery(dbs.Dialect, sql), args)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	result.Data = payload.FieldId
	return result
}
