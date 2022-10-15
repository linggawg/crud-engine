package usecases_test

import (
	"context"
	"database/sql"
	modelsDbs "engine/bin/modules/dbs/models/domain"
	"engine/bin/modules/engine/helpers"
	models "engine/bin/modules/engine/models/domain"
	"engine/bin/modules/engine/models/mocks"
	"engine/bin/modules/engine/usecases"
	modelsQueries "engine/bin/modules/queries/models/domain"
	queriesMocks "engine/bin/modules/queries/models/mocks"
	resourcesMappingModels "engine/bin/modules/resources-mapping/models/domain"
	modelsServices "engine/bin/modules/services/models/domain"
	servicesMocks "engine/bin/modules/services/models/mocks"
	connectionMocks "engine/bin/pkg/databases/connection/mocks"
	"engine/bin/pkg/token"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestEngineQueryUsecase(t *testing.T) {
	suite.Run(t, new(EngineQueryUsecaseTest))
}

type EngineQueryUsecaseTest struct {
	suite.Suite
	servicesQuery *servicesMocks.ServicesPostgreQuery
	queriesQuery  *queriesMocks.QueriesPostgreQuery
	db            *connectionMocks.BulkConnectionPkg
	repo          *mocks.Repository
	usecase       *usecases.EngineQueryUsecase
}

func (s *EngineQueryUsecaseTest) SetupTest() {
	s.servicesQuery = new(servicesMocks.ServicesPostgreQuery)
	s.queriesQuery = new(queriesMocks.QueriesPostgreQuery)
	s.db = new(connectionMocks.BulkConnectionPkg)
	initBulk := *helpers.InitBulkRepository()
	repo, ok := initBulk.GetBulkRepository("mocks").(*mocks.Repository)
	s.Equal(ok, true)
	s.repo = repo
	s.usecase = usecases.NewQueryUsecase(s.queriesQuery, s.servicesQuery, s.db, initBulk)
}

func (s *EngineQueryUsecaseTest) TestGet() {
	table := "users"
	tx := &sqlx.DB{}
	primaryKey := &models.PrimaryKey{Column: "id", Format: "int"}
	dbs := modelsDbs.Dbs{ID: "09582c5e-e7fd-49df-bfa4-a9428ef3a1b4", Name: "db_config", Host: "localhost", Port: 5432, Username: "root", Password: nil, Dialect: "mocks"}
	engineConfig := models.EngineConfig{
		Dbs: dbs,
		ResourcesMappingList: resourcesMappingModels.ResourcesMappingList{resourcesMappingModels.ResourcesMapping{
			ID:           "1",
			ServiceId:    "1",
			SourceOrigin: "id",
			SourceAlias:  "usersId",
		}},
	}
	var result []map[string]interface{}
	result = append(result, map[string]interface{}{"id": "1", "name": "data1"}, map[string]interface{}{"id": "2", "name": "data3"})
	page := 0
	size := 10
	payload := models.GetList{Page: &page, Size: &size, Sort: "id ASC", IsDistinct: true, Filter: "id = '1'", Columns: "", Key: "", Opts: token.Claim{}}
	queries := &modelsQueries.Queries{ID: "1", Key: "get_list_users", QueryDefinition: "SELECT u.id as id, r.name AS role_name, u.username as username, u.created_at as creted_at FROM users u JOIN roles r ON u.role_id = r.id "}
	services := &modelsServices.Services{ID: "1", DbID: dbs.ID, QueryID: &queries.ID, Method: "GET", ServiceUrl: &table}
	s.Run("success-get", func() {
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, table).Return(primaryKey, nil)
		s.repo.On("CountData", mock.Anything, tx, mock.AnythingOfType("string")).Return(int64(2), nil)
		s.repo.On("FindData", mock.Anything, tx, mock.AnythingOfType("string")).Return(result, nil)
		res := s.usecase.Get(context.TODO(), engineConfig, table, &payload)
		s.Nil(res.Error)
		s.queriesQuery.AssertNotCalled(s.T(), "FindOneByKey")
	})
	s.Run("success-get-by-key", func() {
		s.SetupTest()
		payload.Key = "get_list_users"
		payload.Size = nil
		payload.Page = nil
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.queriesQuery.On("FindOneByKey", mock.Anything, payload.Key).Return(queries, nil)
		s.servicesQuery.On("FindOneByQueryID", mock.Anything, queries.ID).Return(services, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, *services.ServiceUrl).Return(primaryKey, nil)
		s.repo.On("CountData", mock.Anything, tx, mock.AnythingOfType("string")).Return(int64(10), nil)
		s.repo.On("FindData", mock.Anything, tx, mock.AnythingOfType("string")).Return(nil, nil)
		res := s.usecase.Get(context.TODO(), engineConfig, table, &payload)
		s.Nil(res.Error)
	})
	s.Run("error-find-data", func() {
		s.SetupTest()
		payload.Sort = ""
		payload.Size = &size
		payload.Page = &page
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.queriesQuery.On("FindOneByKey", mock.Anything, payload.Key).Return(queries, nil)
		s.servicesQuery.On("FindOneByQueryID", mock.Anything, queries.ID).Return(services, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, *services.ServiceUrl).Return(primaryKey, nil)
		s.repo.On("CountData", mock.Anything, tx, mock.AnythingOfType("string")).Return(int64(10), nil)
		s.repo.On("FindData", mock.Anything, tx, mock.AnythingOfType("string")).Return(nil, sql.ErrTxDone)
		res := s.usecase.Get(context.TODO(), engineConfig, table, &payload)
		s.NotNil(res.Error)
	})
	s.Run("error-count-data", func() {
		s.SetupTest()
		payload.Filter = ""
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.queriesQuery.On("FindOneByKey", mock.Anything, payload.Key).Return(queries, nil)
		s.servicesQuery.On("FindOneByQueryID", mock.Anything, queries.ID).Return(services, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, *services.ServiceUrl).Return(primaryKey, nil)
		s.repo.On("CountData", mock.Anything, tx, mock.AnythingOfType("string")).Return(int64(0), sql.ErrTxDone)
		res := s.usecase.Get(context.TODO(), engineConfig, table, &payload)
		s.NotNil(res.Error)
		s.repo.AssertNotCalled(s.T(), "FindData")
	})
	s.Run("error-get-primary-key", func() {
		s.SetupTest()
		payload.Key = ""
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.repo.On("FindPrimaryKey", mock.Anything, tx, table).Return(nil, sql.ErrNoRows)
		res := s.usecase.Get(context.TODO(), engineConfig, table, &payload)
		s.NotNil(res.Error)
		s.repo.AssertNotCalled(s.T(), "FindData")
	})
	s.Run("error-by-key-queries", func() {
		s.SetupTest()
		payload.Key = "get_list_users"
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, nil)
		s.queriesQuery.On("FindOneByKey", mock.Anything, payload.Key).Return(nil, sql.ErrNoRows)
		res := s.usecase.Get(context.TODO(), engineConfig, table, &payload)
		s.NotNil(res.Error)
		s.repo.AssertNotCalled(s.T(), "FindData")
	})
	s.Run("error-connection", func() {
		s.SetupTest()
		s.db.On("GetBulkConnectionSql", dbs).Return(tx, errors.New("cannot add connection to bulk manager connection"))
		res := s.usecase.Get(context.TODO(), engineConfig, table, &payload)
		s.NotNil(res.Error)
		s.repo.AssertNotCalled(s.T(), "FindData")
	})

}
