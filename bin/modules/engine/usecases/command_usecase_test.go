package usecases_test

import (
	"context"
	"database/sql"
	"engine/bin/config"
	modelsDbs "engine/bin/modules/dbs/models/domain"
	"engine/bin/modules/engine/helpers"
	models "engine/bin/modules/engine/models/domain"
	"engine/bin/modules/engine/models/mocks"
	"engine/bin/modules/engine/usecases"
	resourcesMappingModels "engine/bin/modules/resources-mapping/models/domain"
	connectionMocks "engine/bin/pkg/databases/connection/mocks"
	"engine/bin/pkg/token"
	"engine/bin/pkg/utils"
	"errors"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestEngineCommandUsecase(t *testing.T) {
	suite.Run(t, new(EngineCommandUsecaseTest))
}

type EngineCommandUsecaseTest struct {
	suite.Suite
	db       *connectionMocks.BulkConnectionPkg
	repo     *mocks.Repository
	usecases *usecases.EngineCommandUsecase
}

func (s *EngineCommandUsecaseTest) SetupTest() {
	s.db = new(connectionMocks.BulkConnectionPkg)
	initBulk := *helpers.InitBulkRepository()
	repo, ok := initBulk.GetBulkRepository("mocks").(*mocks.Repository)
	s.Equal(ok, true)
	s.repo = repo
	s.usecases = usecases.NewCommandUsecase(s.db, initBulk)
}
func (s *EngineCommandUsecaseTest) TestSetupDataset() {
	tx := &sqlx.DB{}
	config.GlobalEnv.EnginePassword = ""
	var data []map[string]interface{}
	data = append(data, map[string]interface{}{"id": "1", "name": "data1"}, map[string]interface{}{"id": "2", "name": "data1"})
	sqlRole := `INSERT INTO roles (id, name, created_at, created_by) SELECT $1, $2, $3, $4 WHERE NOT EXISTS ( SELECT name FROM roles WHERE name ILIKE $2 );`
	findRole := fmt.Sprintf("SELECT id FROM roles WHERE name = '%s'", config.GlobalEnv.EngineRole)
	sqlUsers := `INSERT INTO users (id, role_id, username, password, created_at, created_by) SELECT $1, $2, CAST($3 AS VARCHAR), $4, $5, $6 WHERE NOT EXISTS ( SELECT username FROM users WHERE username = $3 );`
	sqlApps := `INSERT INTO apps (id, name, created_at, created_by) SELECT $1, CAST($2 AS VARCHAR), $3, $4 WHERE NOT EXISTS ( SELECT name FROM apps WHERE name ILIKE $2 );`
	findApps := fmt.Sprintf("SELECT id FROM apps WHERE name = '%s'", config.ProjectDirName)
	sqlDbs := `INSERT INTO dbs (id, app_id, name, host, port, username, password, dialect, created_at, created_by)
				SELECT $1, $2, CAST($3 AS VARCHAR), CAST($4 AS VARCHAR), $5, CAST($6 AS VARCHAR), CAST($7 AS VARCHAR), CAST($8 AS VARCHAR), $9, $10 
				WHERE NOT EXISTS ( SELECT id FROM dbs WHERE name = $3 AND host = $4 AND port = $5 AND username = $6 AND password = $7 AND dialect = $8 );`
	s.Run("success", func() {
		s.repo.On("InsertOne", mock.Anything, tx, mock.AnythingOfType("string"), mock.Anything).Return(nil)
		s.repo.On("FindData", mock.Anything, tx, mock.AnythingOfType("string")).Return(data, nil)
		res := s.usecases.SetupDataset(context.TODO(), utils.DialectMocks, tx)
		s.Nil(res.Error)
		s.Equal(res.Data, true)
	})
	s.Run("error-insert-role", func() {
		s.SetupTest()
		s.repo.On("InsertOne", mock.Anything, tx, sqlRole, mock.Anything).Return(sql.ErrTxDone)
		res := s.usecases.SetupDataset(context.TODO(), utils.DialectMocks, tx)
		s.NotNil(res.Error)
	})
	s.Run("error-find-role", func() {
		s.SetupTest()
		s.repo.On("InsertOne", mock.Anything, tx, sqlRole, mock.Anything).Return(nil)
		s.repo.On("FindData", mock.Anything, tx, findRole).Return(nil, sql.ErrNoRows)
		res := s.usecases.SetupDataset(context.TODO(), utils.DialectMocks, tx)
		s.NotNil(res.Error)
	})
	s.Run("error-insert-users", func() {
		s.SetupTest()
		s.repo.On("InsertOne", mock.Anything, tx, sqlRole, mock.Anything).Return(nil)
		s.repo.On("FindData", mock.Anything, tx, findRole).Return(data, nil)
		s.repo.On("InsertOne", mock.Anything, tx, sqlUsers, mock.Anything).Return(sql.ErrTxDone)
		res := s.usecases.SetupDataset(context.TODO(), utils.DialectMocks, tx)
		s.NotNil(res.Error)
	})
	s.Run("error-insert-apps", func() {
		s.SetupTest()
		s.repo.On("InsertOne", mock.Anything, tx, sqlRole, mock.Anything).Return(nil)
		s.repo.On("FindData", mock.Anything, tx, findRole).Return(data, nil)
		s.repo.On("InsertOne", mock.Anything, tx, sqlUsers, mock.Anything).Return(nil)
		s.repo.On("InsertOne", mock.Anything, tx, sqlApps, mock.Anything).Return(sql.ErrTxDone)
		res := s.usecases.SetupDataset(context.TODO(), utils.DialectMocks, tx)
		s.NotNil(res.Error)
	})
	s.Run("error-find-apps", func() {
		s.SetupTest()
		s.repo.On("InsertOne", mock.Anything, tx, sqlRole, mock.Anything).Return(nil)
		s.repo.On("FindData", mock.Anything, tx, findRole).Return(data, nil)
		s.repo.On("InsertOne", mock.Anything, tx, sqlUsers, mock.Anything).Return(nil)
		s.repo.On("InsertOne", mock.Anything, tx, sqlApps, mock.Anything).Return(nil)
		s.repo.On("FindData", mock.Anything, tx, findApps).Return(nil, sql.ErrNoRows)
		res := s.usecases.SetupDataset(context.TODO(), utils.DialectMocks, tx)
		s.NotNil(res.Error)
	})
	s.Run("error-insert-dbs", func() {
		s.SetupTest()
		s.repo.On("InsertOne", mock.Anything, tx, sqlRole, mock.Anything).Return(nil)
		s.repo.On("FindData", mock.Anything, tx, findRole).Return(data, nil)
		s.repo.On("InsertOne", mock.Anything, tx, sqlUsers, mock.Anything).Return(nil)
		s.repo.On("InsertOne", mock.Anything, tx, sqlApps, mock.Anything).Return(nil)
		s.repo.On("FindData", mock.Anything, tx, findApps).Return(data, nil)
		s.repo.On("InsertOne", mock.Anything, tx, sqlDbs, mock.Anything).Return(sql.ErrTxDone)
		res := s.usecases.SetupDataset(context.TODO(), utils.DialectMocks, tx)
		s.NotNil(res.Error)
	})
}
func (s *EngineCommandUsecaseTest) TestInsert() {
	tx := &sqlx.DB{}
	dbs := modelsDbs.Dbs{ID: "09582c5e-e7fd-49df-bfa4-a9428ef3a1b4", Name: "db_config", Host: "localhost", Port: 5432, Username: "root", Password: nil, Dialect: "mocks"}
	engineConfig := models.EngineConfig{
		Dbs:                  dbs,
		ResourcesMappingList: resourcesMappingModels.ResourcesMappingList{resourcesMappingModels.ResourcesMapping{ID: "1", ServiceId: "1", SourceOrigin: "id", SourceAlias: "usersId"}, resourcesMappingModels.ResourcesMapping{ID: "2", ServiceId: "1", SourceOrigin: "name", SourceAlias: "usersName"}, resourcesMappingModels.ResourcesMapping{ID: "3", ServiceId: "1", SourceOrigin: "created_by", SourceAlias: "createdBy"}},
	}
	payload := &models.EngineRequest{Table: "users", FieldId: "id", Value: "1", Data: map[string]interface{}{"usersId": "1", "usersName": "data1"}, Opts: token.Claim{}}
	primaryKey := &models.PrimaryKey{Column: "id", Format: "character varying(255)"}
	var schema []models.InformationSchema
	schema = append(schema, models.InformationSchema{ColumName: "id", IsNullable: "NO"}, models.InformationSchema{ColumName: "name", IsNullable: "YES"})
	s.Run("success", func() {
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, payload.Table).Return(primaryKey, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(schema, nil)
		s.repo.On("InsertOne", mock.Anything, tx, mock.AnythingOfType("string"), mock.Anything).Return(nil)
		res := s.usecases.Insert(context.TODO(), engineConfig, payload)
		s.Nil(res.Error)
	})
	s.Run("error-connection", func() {
		s.SetupTest()
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, errors.New("cannot add connection to bulk manager connection"))
		res := s.usecases.Insert(context.TODO(), engineConfig, payload)
		s.NotNil(res.Error)
	})
	s.Run("error-get-primary-key", func() {
		s.SetupTest()
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, payload.Table).Return(nil, sql.ErrNoRows)
		res := s.usecases.Insert(context.TODO(), engineConfig, payload)
		s.NotNil(res.Error)
	})
	s.Run("error-get-schema", func() {
		s.SetupTest()
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, payload.Table).Return(primaryKey, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(nil, sql.ErrTxDone)
		res := s.usecases.Insert(context.TODO(), engineConfig, payload)
		s.NotNil(res.Error)
	})
	s.Run("error-if-value-body-nil", func() {
		s.SetupTest()
		schema = append(schema, models.InformationSchema{
			ColumName:  "created_by",
			IsNullable: "NO",
		})
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, payload.Table).Return(primaryKey, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(schema, nil)
		res := s.usecases.Insert(context.TODO(), engineConfig, payload)
		s.NotNil(res.Error)
		s.repo.AssertNotCalled(s.T(), "InsertOne")
	})
	s.Run("error-if-body-cannot-be-empty", func() {
		s.SetupTest()
		payload.Data = map[string]interface{}{"id": "1", "usersId": "userName", "createdBy": nil}
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, payload.Table).Return(primaryKey, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(schema, nil)
		res := s.usecases.Insert(context.TODO(), engineConfig, payload)
		s.NotNil(res.Error)
		s.repo.AssertNotCalled(s.T(), "InsertOne")
	})
	s.Run("error-insert", func() {
		s.SetupTest()
		payload.Data = map[string]interface{}{"usersId": "1", "name": "userName", "createdBy": "1"}
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, payload.Table).Return(primaryKey, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(schema, nil)
		s.repo.On("InsertOne", mock.Anything, tx, mock.AnythingOfType("string"), mock.Anything).Return(sql.ErrTxDone)
		res := s.usecases.Insert(context.TODO(), engineConfig, payload)
		s.NotNil(res.Error)
	})
}

func (s *EngineCommandUsecaseTest) TestUpdate() {
	tx := &sqlx.DB{}
	dbs := modelsDbs.Dbs{ID: "09582c5e-e7fd-49df-bfa4-a9428ef3a1b4", Name: "db_config", Host: "localhost", Port: 5432, Username: "root", Password: nil, Dialect: "mocks"}
	engineConfig := models.EngineConfig{
		Dbs:                  dbs,
		ResourcesMappingList: resourcesMappingModels.ResourcesMappingList{resourcesMappingModels.ResourcesMapping{ID: "1", ServiceId: "1", SourceOrigin: "id", SourceAlias: "usersId"}, resourcesMappingModels.ResourcesMapping{ID: "2", ServiceId: "1", SourceOrigin: "name", SourceAlias: "usersName"}, resourcesMappingModels.ResourcesMapping{ID: "3", ServiceId: "1", SourceOrigin: "created_by", SourceAlias: "createdBy"}},
	}
	payload := &models.EngineRequest{Table: "users", FieldId: "usersId", Value: "1", Data: map[string]interface{}{"usersId": "1", "usersName": "data1"}, Opts: token.Claim{}}
	primaryKey := &models.PrimaryKey{Column: "id", Format: "character varying(255)"}
	var schema []models.InformationSchema
	schema = append(schema, models.InformationSchema{ColumName: "id", IsNullable: "NO"}, models.InformationSchema{ColumName: "name", IsNullable: "YES"})
	s.Run("success", func() {
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(schema, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, payload.Table).Return(primaryKey, nil)
		s.repo.On("UpdateOne", mock.Anything, tx, mock.AnythingOfType("string"), mock.Anything).Return(nil)
		s.Nil(s.usecases.Update(context.TODO(), engineConfig, payload).Error)
		s.Nil(s.usecases.Patch(context.TODO(), engineConfig, payload).Error)
	})
	s.Run("error-connection", func() {
		s.SetupTest()
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, errors.New("cannot add connection to bulk manager connection"))
		s.NotNil(s.usecases.Update(context.TODO(), engineConfig, payload).Error)
		s.NotNil(s.usecases.Patch(context.TODO(), engineConfig, payload).Error)
		s.repo.AssertNotCalled(s.T(), "UpdateOne")
	})
	s.Run("error-get-schema", func() {
		s.SetupTest()
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(nil, sql.ErrTxDone)
		s.NotNil(s.usecases.Update(context.TODO(), engineConfig, payload).Error)
		s.NotNil(s.usecases.Patch(context.TODO(), engineConfig, payload).Error)
		s.repo.AssertNotCalled(s.T(), "UpdateOne")
	})
	s.Run("error-get-primary-key", func() {
		s.SetupTest()
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(schema, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, payload.Table).Return(nil, sql.ErrNoRows)
		s.NotNil(s.usecases.Update(context.TODO(), engineConfig, payload).Error)
		s.NotNil(s.usecases.Patch(context.TODO(), engineConfig, payload).Error)
		s.repo.AssertNotCalled(s.T(), "UpdateOne")
	})
	s.Run("error-update", func() {
		s.SetupTest()
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(schema, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, payload.Table).Return(primaryKey, nil)
		s.repo.On("UpdateOne", mock.Anything, tx, mock.AnythingOfType("string"), mock.Anything).Return(sql.ErrTxDone)
		s.NotNil(s.usecases.Update(context.TODO(), engineConfig, payload).Error)
		s.NotNil(s.usecases.Patch(context.TODO(), engineConfig, payload).Error)
		s.repo.AssertNotCalled(s.T(), "UpdateOne")
	})
	s.Run("error-if-body-cannot-be-empty", func() {
		s.SetupTest()
		schema = append(schema, models.InformationSchema{ColumName: "created_by", IsNullable: "NO"})
		payload.Data = map[string]interface{}{"usersId": "1", "usersName": "data1", "createdBy": nil}
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(schema, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, payload.Table).Return(primaryKey, nil)
		s.NotNil(s.usecases.Update(context.TODO(), engineConfig, payload).Error)
		s.NotNil(s.usecases.Patch(context.TODO(), engineConfig, payload).Error)
		s.repo.AssertNotCalled(s.T(), "UpdateOne")
	})
	s.Run("error-field-id-not-found", func() {
		s.SetupTest()
		payload.FieldId = "users_id"
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(schema, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, payload.Table).Return(primaryKey, nil)
		s.NotNil(s.usecases.Update(context.TODO(), engineConfig, payload).Error)
		s.NotNil(s.usecases.Patch(context.TODO(), engineConfig, payload).Error)
		s.repo.AssertNotCalled(s.T(), "UpdateOne")
	})
}

func (s *EngineCommandUsecaseTest) TestDelete() {
	tx := &sqlx.DB{}
	dbs := modelsDbs.Dbs{ID: "09582c5e-e7fd-49df-bfa4-a9428ef3a1b4", Name: "db_config", Host: "localhost", Port: 5432, Username: "root", Password: nil, Dialect: "mocks"}
	engineConfig := models.EngineConfig{
		Dbs:                  dbs,
		ResourcesMappingList: resourcesMappingModels.ResourcesMappingList{resourcesMappingModels.ResourcesMapping{ID: "1", ServiceId: "1", SourceOrigin: "id", SourceAlias: "usersId"}},
	}
	payload := &models.EngineRequest{Table: "users", FieldId: "usersId", Value: "1", Opts: token.Claim{}}
	var schema []models.InformationSchema
	schema = append(schema, models.InformationSchema{ColumName: "id", IsNullable: "NO"}, models.InformationSchema{ColumName: "name", IsNullable: "YES"})
	s.Run("success", func() {
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(schema, nil)
		s.repo.On("DeleteOne", mock.Anything, tx, mock.AnythingOfType("string"), mock.Anything).Return(nil)
		s.Nil(s.usecases.Delete(context.TODO(), engineConfig, payload).Error)
	})
	s.Run("error-connection", func() {
		s.SetupTest()
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, errors.New("cannot add connection to bulk manager connection"))
		s.NotNil(s.usecases.Delete(context.TODO(), engineConfig, payload).Error)
		s.repo.AssertNotCalled(s.T(), "DeleteOne")
	})
	s.Run("error-get-schema", func() {
		s.SetupTest()
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(nil, sql.ErrTxDone)
		s.NotNil(s.usecases.Delete(context.TODO(), engineConfig, payload).Error)
		s.repo.AssertNotCalled(s.T(), "DeleteOne")
	})
	s.Run("error-delete", func() {
		s.SetupTest()
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(schema, nil)
		s.repo.On("DeleteOne", mock.Anything, tx, mock.AnythingOfType("string"), mock.Anything).Return(sql.ErrTxDone)
		s.NotNil(s.usecases.Delete(context.TODO(), engineConfig, payload).Error)
	})
	s.Run("error-field-id-not-found", func() {
		s.SetupTest()
		payload.FieldId = "users_id"
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("SelectInformationSchema", mock.Anything, tx, payload.Table).Return(schema, nil)
		s.NotNil(s.usecases.Delete(context.TODO(), engineConfig, payload).Error)
		s.repo.AssertNotCalled(s.T(), "DeleteOne")
	})
}
