package usecases

import (
	"context"
	"engine/bin/config"
	"engine/bin/modules/engine/helpers"
	models "engine/bin/modules/engine/models/domain"
	resourcesMappingModels "engine/bin/modules/resources-mapping/models/domain"
	"engine/bin/pkg/databases/connection"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/utils"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"
)

type EngineCommandUsecase struct {
	db   connection.Connection
	repo *helpers.BulkRepository
}

func NewCommandUsecase(db connection.Connection, repo helpers.BulkRepository) *EngineCommandUsecase {
	return &EngineCommandUsecase{
		db:   db,
		repo: &repo,
	}
}

func (h *EngineCommandUsecase) SetupDataset(ctx context.Context, dialect string, db *sqlx.DB) (result utils.Result) {
	repository := h.repo.GetBulkRepository(dialect)

	//insert default app
	var apps []interface{}
	sql := `INSERT INTO apps (id, name, created_at, created_by) SELECT $1, CAST($2 AS VARCHAR), $3, $4 WHERE NOT EXISTS (SELECT name FROM apps WHERE LOWER(name) LIKE LOWER($2));`
	apps = append(apps, uuid.New().String(), config.ProjectDirName, time.Now(), config.ProjectDirName)
	err := repository.InsertOne(ctx, db, sql, apps)
	if err != nil {
		errObj := httpError.NewConflict()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	//insert default dbs
	res, err := repository.FindData(ctx, db, fmt.Sprintf("SELECT id FROM apps WHERE name = '%s'", config.ProjectDirName))
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	if len(res) > 0 {
		appId := res[0]["id"]
		var dbs []interface{}
		sql = `INSERT INTO dbs (id, app_id, name, host, port, username, password, dialect, created_at, created_by)
				SELECT $1, $2, CAST($3 AS VARCHAR), CAST($4 AS VARCHAR), $5, CAST($6 AS VARCHAR), CAST($7 AS VARCHAR), CAST($8 AS VARCHAR), $9, $10 
				WHERE NOT EXISTS ( SELECT id FROM dbs WHERE name = $3 AND host = $4 AND port = $5 AND username = $6 AND password = $7 AND dialect = $8 );`
		dbs = append(dbs, uuid.New().String(), appId, config.GlobalEnv.DBName, config.GlobalEnv.DBHost, config.GlobalEnv.DBPort, config.GlobalEnv.DBUser, config.GlobalEnv.DBPassword, config.GlobalEnv.DBDialect, time.Now(), config.ProjectDirName)
		err = repository.InsertOne(ctx, db, sql, dbs)
		if err != nil {
			errObj := httpError.NewConflict()
			errObj.Message = err.Error()
			result.Error = errObj
			return result
		}
	}

	//insert default role
	var role []interface{}
	sql = `INSERT INTO roles (id, name, created_at, created_by) SELECT $1, CAST($2 AS VARCHAR), $3, $4 WHERE NOT EXISTS (SELECT name FROM roles WHERE LOWER(name) LIKE LOWER($2));`
	role = append(role, uuid.New().String(), utils.RoleAdmin, time.Now(), utils.RoleAdmin)
	err = repository.InsertOne(ctx, db, sql, role)
	if err != nil {
		errObj := httpError.NewConflict()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	//insert default dbs
	res, err = repository.FindData(ctx, db, fmt.Sprintf("SELECT id FROM roles WHERE name = '%s'", utils.RoleAdmin))
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	if len(res) > 0 {
		roleId := res[0]["id"]
		var user []interface{}
		sql = `INSERT INTO users (id, role_id, username, password, created_at, created_by) SELECT $1, $2, CAST($3 AS VARCHAR), CAST($4 AS VARCHAR), $5, $6 WHERE NOT EXISTS (SELECT name FROM users WHERE LOWER(name) LIKE LOWER($3));`
		user = append(user, uuid.New().String(), roleId, utils.RoleAdmin, utils.RoleAdmin, time.Now(), utils.RoleAdmin)
		err = repository.InsertOne(ctx, db, sql, user)
		if err != nil {
			errObj := httpError.NewConflict()
			errObj.Message = err.Error()
			result.Error = errObj
			return result
		}
	}

	result.Data = true
	return result
}

func (h *EngineCommandUsecase) Insert(ctx context.Context, engineConfig models.EngineConfig, payload *models.EngineRequest) (result utils.Result) {
	repository := h.repo.GetBulkRepository(engineConfig.Dbs.Dialect)
	tx, err := h.db.GetBulkConnectionSql(engineConfig.Dbs)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	primaryKey, err := repository.FindPrimaryKey(ctx, tx, payload.Table)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Undefined primary key"
		result.Error = errObj
		return result
	}
	schema, err := repository.SelectInformationSchema(ctx, tx, payload.Table)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	sql, args, err := h.SqlStatementInsert(engineConfig, primaryKey, payload, schema)
	if err != nil {
		errObj := httpError.NewBadRequest()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	err = repository.InsertOne(ctx, tx, sql, args)
	if err != nil {
		errObj := httpError.NewConflict()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	result.Data = payload.Data
	return result
}

func (h *EngineCommandUsecase) SqlStatementInsert(engineConfig models.EngineConfig, primaryKey *models.PrimaryKey, payload *models.EngineRequest, informationSchemas []models.InformationSchema) (sql string, args []interface{}, err error) {
	var listColumnReq = payload.Data
	if len(engineConfig.ResourcesMappingList) > 0 {
		listColumnReq = helpers.ConvertToSourceOrigin(payload.Data, engineConfig.ResourcesMappingList)
	}

	var columns, values []string
	for key := range listColumnReq {
		if key != primaryKey.Column {
			columns = append(columns, key)
			values = append(values, "?")
			if listColumnReq[key] == nil {
				for _, i := range informationSchemas {
					if i.IsNullable == "NO" && i.ColumName == key {
						return "", args, ErrorMessageMapping(i.ColumName, engineConfig.ResourcesMappingList)
					}
				}
			}
			args = append(args, listColumnReq[key])
		}
	}
	if !strings.EqualFold(primaryKey.Format, "int") {
		columns = append(columns, primaryKey.Column)
		values = append(values, "?")
		args = append(args, uuid.New().String())
	}

	for _, i := range informationSchemas {
		if i.IsNullable == "NO" && i.ColumName != primaryKey.Column {
			if !FindValueInListString(columns, i.ColumName) {
				return "", args, ErrorMessageMapping(i.ColumName, engineConfig.ResourcesMappingList)
			}
		}
	}
	for i := range columns {
		columns[i] = fmt.Sprintf("`" + columns[i] + "`")
	}
	sql = fmt.Sprintf(utils.QueryInsert, payload.Table, strings.Join(columns, " ,"), strings.Join(values, ","))
	return helpers.SetQuery(engineConfig.Dbs.Dialect, sql), args, nil
}

func (h *EngineCommandUsecase) Update(ctx context.Context, engineConfig models.EngineConfig, payload *models.EngineRequest) (result utils.Result) {
	repository := h.repo.GetBulkRepository(engineConfig.Dbs.Dialect)

	tx, err := h.db.GetBulkConnectionSql(engineConfig.Dbs)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	schema, err := repository.SelectInformationSchema(ctx, tx, payload.Table)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	if !CheckFieldId(engineConfig, schema, payload) {
		errObj := httpError.NewBadRequest()
		errObj.Message = fmt.Sprintf("field_id '%s' is not found", payload.FieldId)
		result.Error = errObj
		return result
	}

	primaryKey, err := repository.FindPrimaryKey(ctx, tx, payload.Table)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Undefined primary key"
		result.Error = errObj
		return result
	}

	sql, args, err := h.SqlStatementUpdate(engineConfig, payload, primaryKey, schema)
	if err != nil {
		errObj := httpError.NewBadRequest()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	err = repository.UpdateOne(ctx, tx, sql, args)
	if err != nil {
		errObj := httpError.NewConflict()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	result.Data = payload.Data
	return result
}

func (h *EngineCommandUsecase) Patch(ctx context.Context, engineConfig models.EngineConfig, payload *models.EngineRequest) (result utils.Result) {
	repository := h.repo.GetBulkRepository(engineConfig.Dbs.Dialect)

	tx, err := h.db.GetBulkConnectionSql(engineConfig.Dbs)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	schema, err := repository.SelectInformationSchema(ctx, tx, payload.Table)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	if !CheckFieldId(engineConfig, schema, payload) {
		errObj := httpError.NewBadRequest()
		errObj.Message = fmt.Sprintf("field_id '%s' is not found", payload.FieldId)
		result.Error = errObj
		return result
	}

	primaryKey, err := repository.FindPrimaryKey(ctx, tx, payload.Table)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Undefined primary key"
		result.Error = errObj
		return result
	}

	sql, args, err := h.SqlStatementPatch(engineConfig, payload, primaryKey, schema)
	if err != nil {
		errObj := httpError.NewBadRequest()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	err = repository.UpdateOne(ctx, tx, sql, args)
	if err != nil {
		errObj := httpError.NewConflict()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	result.Data = payload.Data
	return result
}

func (h *EngineCommandUsecase) SqlStatementUpdate(engineConfig models.EngineConfig, payload *models.EngineRequest, primaryKey *models.PrimaryKey, informationSchemas []models.InformationSchema) (sql string, args []interface{}, err error) {
	var (
		setData       []string
		listColumnReq = payload.Data
	)
	if len(engineConfig.ResourcesMappingList) > 0 {
		listColumnReq = helpers.ConvertToSourceOrigin(payload.Data, engineConfig.ResourcesMappingList)
	}
	for _, i := range informationSchemas {
		if i.ColumName != primaryKey.Column {
			if i.IsNullable == "NO" && listColumnReq[i.ColumName] == nil {
				return "", args, ErrorMessageMapping(i.ColumName, engineConfig.ResourcesMappingList)
			}
			setData = append(setData, fmt.Sprintf("`%s`=?", i.ColumName))
			args = append(args, listColumnReq[i.ColumName])
		}
	}

	args = append(args, payload.Value)
	sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", payload.Table, strings.Join(setData, ","), payload.FieldId)
	return helpers.SetQuery(engineConfig.Dbs.Dialect, sql), args, nil
}

func (h *EngineCommandUsecase) SqlStatementPatch(engineConfig models.EngineConfig, payload *models.EngineRequest, primaryKey *models.PrimaryKey, informationSchemas []models.InformationSchema) (sql string, args []interface{}, err error) {
	var (
		setData       []string
		listColumnReq = payload.Data
	)
	if len(engineConfig.ResourcesMappingList) > 0 {
		listColumnReq = helpers.ConvertToSourceOrigin(payload.Data, engineConfig.ResourcesMappingList)
	}
	for key := range listColumnReq {
		if key != primaryKey.Column {
			if listColumnReq[key] == nil {
				for _, i := range informationSchemas {
					if i.ColumName == key && i.IsNullable == "NO" {
						return "", args, ErrorMessageMapping(i.ColumName, engineConfig.ResourcesMappingList)
					}
				}
			}
			setData = append(setData, fmt.Sprintf("`%s`=?", key))
			args = append(args, listColumnReq[key])
		}
	}
	args = append(args, payload.Value)
	sql = fmt.Sprintf(utils.QueryUpdate, payload.Table, strings.Join(setData, ","), payload.FieldId)
	return helpers.SetQuery(engineConfig.Dbs.Dialect, sql), args, nil
}

func (h *EngineCommandUsecase) Delete(ctx context.Context, engineConfig models.EngineConfig, payload *models.EngineRequest) (result utils.Result) {
	var args []interface{}
	repository := h.repo.GetBulkRepository(engineConfig.Dbs.Dialect)

	tx, err := h.db.GetBulkConnectionSql(engineConfig.Dbs)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}

	schema, err := repository.SelectInformationSchema(ctx, tx, payload.Table)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	if !CheckFieldId(engineConfig, schema, payload) {
		errObj := httpError.NewBadRequest()
		errObj.Message = fmt.Sprintf("field_id '%s' is not found", payload.FieldId)
		result.Error = errObj
		return result
	}
	sql := fmt.Sprintf(utils.QueryDelete, payload.Table, payload.FieldId)
	args = append(args, payload.Value)
	err = repository.DeleteOne(ctx, tx, helpers.SetQuery(engineConfig.Dbs.Dialect, sql), args)
	if err != nil {
		errObj := httpError.NewConflict()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	result.Data = payload.FieldId
	return result
}

func FindValueInListString(s []string, str string) bool {
	for i := range s {
		if s[i] == str {
			return true
		}
	}
	return false
}

func ErrorMessageMapping(columnName string, list resourcesMappingModels.ResourcesMappingList) error {
	validateColumn := columnName
	for _, rm := range list {
		if rm.SourceOrigin == columnName {
			validateColumn = rm.SourceAlias
			break
		}
	}
	return fmt.Errorf("validation for '%s' failed on the 'required' tag", validateColumn)
}

func CheckFieldId(engineConfig models.EngineConfig, schema []models.InformationSchema, payload *models.EngineRequest) bool {
	for i := range engineConfig.ResourcesMappingList {
		if payload.FieldId == engineConfig.ResourcesMappingList[i].SourceAlias {
			payload.FieldId = engineConfig.ResourcesMappingList[i].SourceOrigin
			break
		}
	}
	for _, i := range schema {
		if i.ColumName == payload.FieldId {
			return true
		}
	}
	return false
}
