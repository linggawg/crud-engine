package usecases_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	modelsDbs "engine/bin/modules/dbs/models/domain"
	"engine/bin/modules/dbs/models/mocks"
	"engine/bin/modules/dbs/usecases"
	modelsQueries "engine/bin/modules/queries/models/domain"
	queriesMocks "engine/bin/modules/queries/models/mocks"
	modelsRM "engine/bin/modules/resources-mapping/models/domain"
	resourcesMappingMocks "engine/bin/modules/resources-mapping/models/mocks"
	modelsServices "engine/bin/modules/services/models/domain"
	servicesMocks "engine/bin/modules/services/models/mocks"
	modelsUsersServices "engine/bin/modules/users-services/models/domain"
	usersServicesMocks "engine/bin/modules/users-services/models/mocks"
	usersMocks "engine/bin/modules/users/models/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestDbsQueryUsecase(t *testing.T) {
	suite.Run(t, new(DbsQueryUsecaseTest))
}

type DbsQueryUsecaseTest struct {
	suite.Suite
	dbsPostgreQuery    *mocks.DbsPostgreQuery
	servicesQuery      *servicesMocks.ServicesPostgreQuery
	usersServicesQuery *usersServicesMocks.UsersServicesPostgreQuery
	usersQuery         *usersMocks.UsersPostgreQuery
	queriesQuery       *queriesMocks.QueriesPostgreQuery
	resourcesMapping   *resourcesMappingMocks.ResourcesMappingPostgre
	usecases           *usecases.DbsQueryUsecase
}

func (s *DbsQueryUsecaseTest) SetupTest() {
	s.dbsPostgreQuery = new(mocks.DbsPostgreQuery)
	s.servicesQuery = new(servicesMocks.ServicesPostgreQuery)
	s.usersServicesQuery = new(usersServicesMocks.UsersServicesPostgreQuery)
	s.usersQuery = new(usersMocks.UsersPostgreQuery)
	s.queriesQuery = new(queriesMocks.QueriesPostgreQuery)
	s.resourcesMapping = new(resourcesMappingMocks.ResourcesMappingPostgre)
	s.usecases = usecases.NewQueryUsecase(s.dbsPostgreQuery, s.queriesQuery, s.servicesQuery, s.usersServicesQuery, s.usersQuery, s.resourcesMapping)
}

func (s *DbsQueryUsecaseTest) TestGetDbsConnection() {
	userId := "test1"
	method := "GET"
	serviceUrl := "users"
	services := &modelsServices.Services{ID: "ID1", DbID: "DbID2"}
	queries := &modelsQueries.Queries{ID: "queries1"}
	s.Run("success-find", func() {
		s.servicesQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, serviceUrl, method).Return(services, nil)
		s.usersServicesQuery.On("FindOneByServiceIDAndUserId", mock.Anything, services.ID, userId).Return(&modelsUsersServices.UsersServices{}, nil)
		s.resourcesMapping.On("FindByServiceId", mock.Anything, services.ID).Return(modelsRM.ResourcesMappingList{modelsRM.ResourcesMapping{ID: "1"}}, nil)
		s.dbsPostgreQuery.On("FindOneByID", mock.Anything, services.DbID).Return(&modelsDbs.Dbs{ID: services.DbID}, nil)
		res := s.usecases.GetDbsConnection(context.TODO(), userId, method, serviceUrl, "")
		s.Nil(res.Error)
		s.queriesQuery.AssertNotCalled(s.T(), "FindOneByKey")
		s.servicesQuery.AssertNotCalled(s.T(), "FindOneByQueryID")
	})
	s.Run("success-find-query", func() {
		s.SetupTest()
		key := "get_users"
		s.queriesQuery.On("FindOneByKey", mock.Anything, key).Return(queries, nil)
		s.servicesQuery.On("FindOneByQueryID", mock.Anything, queries.ID).Return(services, nil)
		s.usersServicesQuery.On("FindOneByServiceIDAndUserId", mock.Anything, services.ID, userId).Return(&modelsUsersServices.UsersServices{}, nil)
		s.resourcesMapping.On("FindByServiceId", mock.Anything, services.ID).Return(modelsRM.ResourcesMappingList{modelsRM.ResourcesMapping{ID: "1"}}, nil)
		s.dbsPostgreQuery.On("FindOneByID", mock.Anything, services.DbID).Return(&modelsDbs.Dbs{ID: services.DbID}, nil)
		res := s.usecases.GetDbsConnection(context.TODO(), userId, method, "query", key)
		s.Nil(res.Error)
		s.servicesQuery.AssertNotCalled(s.T(), "FindOneByServiceUrlAndMethodAndQueryIsNull")
	})
	s.Run("error-find-query-services", func() {
		s.SetupTest()
		key := "get_users"
		s.queriesQuery.On("FindOneByKey", mock.Anything, key).Return(queries, nil)
		s.servicesQuery.On("FindOneByQueryID", mock.Anything, queries.ID).Return(nil, sql.ErrNoRows)
		res := s.usecases.GetDbsConnection(context.TODO(), userId, method, "query", key)
		s.NotNil(res.Error)
		s.servicesQuery.AssertNotCalled(s.T(), "FindOneByQueryID")
		s.usersServicesQuery.AssertNotCalled(s.T(), "FindOneByServiceIDAndUserId")
		s.resourcesMapping.AssertNotCalled(s.T(), "FindByServiceId")
		s.dbsPostgreQuery.AssertNotCalled(s.T(), "FindOneByID")
	})
	s.Run("error-find-query", func() {
		s.SetupTest()
		key := "get_users"
		s.queriesQuery.On("FindOneByKey", mock.Anything, key).Return(nil, sql.ErrNoRows)
		res := s.usecases.GetDbsConnection(context.TODO(), userId, method, "query", key)
		s.NotNil(res.Error)
		s.usersServicesQuery.AssertNotCalled(s.T(), "FindOneByServiceIDAndUserId")
		s.resourcesMapping.AssertNotCalled(s.T(), "FindByServiceId")
		s.dbsPostgreQuery.AssertNotCalled(s.T(), "FindOneByID")
	})
	s.Run("error-services-query", func() {
		s.SetupTest()
		s.servicesQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, serviceUrl, method).Return(nil, sql.ErrNoRows)
		res := s.usecases.GetDbsConnection(context.TODO(), userId, method, serviceUrl, "")
		s.NotNil(res.Error)
		s.usersServicesQuery.AssertNotCalled(s.T(), "FindOneByServiceIDAndUserId")
		s.resourcesMapping.AssertNotCalled(s.T(), "FindByServiceId")
		s.dbsPostgreQuery.AssertNotCalled(s.T(), "FindOneByID")
	})
	s.Run("error-users-services-query", func() {
		s.SetupTest()
		s.servicesQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, serviceUrl, method).Return(services, nil)
		s.usersServicesQuery.On("FindOneByServiceIDAndUserId", mock.Anything, services.ID, userId).Return(nil, sql.ErrNoRows)
		res := s.usecases.GetDbsConnection(context.TODO(), userId, method, serviceUrl, "")
		s.NotNil(res.Error)
		s.resourcesMapping.AssertNotCalled(s.T(), "FindByServiceId")
		s.dbsPostgreQuery.AssertNotCalled(s.T(), "FindOneByID")
	})
	s.Run("error-resources-mapping", func() {
		s.SetupTest()
		s.servicesQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, serviceUrl, method).Return(services, nil)
		s.usersServicesQuery.On("FindOneByServiceIDAndUserId", mock.Anything, services.ID, userId).Return(&modelsUsersServices.UsersServices{}, nil)
		s.resourcesMapping.On("FindByServiceId", mock.Anything, services.ID).Return(modelsRM.ResourcesMappingList{}, sql.ErrTxDone)
		res := s.usecases.GetDbsConnection(context.TODO(), userId, method, serviceUrl, "")
		s.NotNil(res.Error)
		s.dbsPostgreQuery.AssertNotCalled(s.T(), "FindOneByID")
	})
	s.Run("error", func() {
		s.SetupTest()
		s.servicesQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, serviceUrl, method).Return(services, nil)
		s.usersServicesQuery.On("FindOneByServiceIDAndUserId", mock.Anything, services.ID, userId).Return(&modelsUsersServices.UsersServices{}, nil)
		s.resourcesMapping.On("FindByServiceId", mock.Anything, services.ID).Return(modelsRM.ResourcesMappingList{modelsRM.ResourcesMapping{ID: "1"}}, nil)
		s.dbsPostgreQuery.On("FindOneByID", mock.Anything, services.DbID).Return(nil, sql.ErrNoRows)
		res := s.usecases.GetDbsConnection(context.TODO(), userId, method, serviceUrl, "")
		s.NotNil(res.Error)
		s.dbsPostgreQuery.AssertExpectations(s.T())
	})
	s.Run("internal-server-error", func() {
		s.SetupTest()
		key := "get_users"
		s.queriesQuery.On("FindOneByKey", mock.Anything, key).Return(nil, driver.ErrBadConn)
		res := s.usecases.GetDbsConnection(context.TODO(), userId, method, "query", key)
		s.NotNil(res.Error)
		s.servicesQuery.AssertNotCalled(s.T(), "FindOneByQueryID")
		s.usersServicesQuery.AssertNotCalled(s.T(), "FindOneByServiceIDAndUserId")
		s.resourcesMapping.AssertNotCalled(s.T(), "FindByServiceId")
		s.dbsPostgreQuery.AssertNotCalled(s.T(), "FindOneByID")
	})
	s.Run("internal-server-error", func() {
		s.SetupTest()
		key := "get_users"
		s.queriesQuery.On("FindOneByKey", mock.Anything, key).Return(queries, nil)
		s.servicesQuery.On("FindOneByQueryID", mock.Anything, queries.ID).Return(nil, driver.ErrBadConn)
		res := s.usecases.GetDbsConnection(context.TODO(), userId, method, "query", key)
		s.NotNil(res.Error)
		s.usersServicesQuery.AssertNotCalled(s.T(), "FindOneByServiceIDAndUserId")
		s.resourcesMapping.AssertNotCalled(s.T(), "FindByServiceId")
		s.dbsPostgreQuery.AssertNotCalled(s.T(), "FindOneByID")
	})
	s.Run("internal-server-error", func() {
		s.SetupTest()
		s.servicesQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, serviceUrl, method).Return(nil, driver.ErrBadConn)
		res := s.usecases.GetDbsConnection(context.TODO(), userId, method, serviceUrl, "")
		s.NotNil(res.Error)
		s.usersServicesQuery.AssertNotCalled(s.T(), "FindOneByServiceIDAndUserId")
		s.resourcesMapping.AssertNotCalled(s.T(), "FindByServiceId")
		s.dbsPostgreQuery.AssertNotCalled(s.T(), "FindOneByID")
	})
	s.Run("internal-server-error", func() {
		s.SetupTest()
		s.servicesQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, serviceUrl, method).Return(services, nil)
		s.usersServicesQuery.On("FindOneByServiceIDAndUserId", mock.Anything, services.ID, userId).Return(nil, driver.ErrBadConn)
		res := s.usecases.GetDbsConnection(context.TODO(), userId, method, serviceUrl, "")
		s.NotNil(res.Error)
		s.resourcesMapping.AssertNotCalled(s.T(), "FindByServiceId")
		s.dbsPostgreQuery.AssertNotCalled(s.T(), "FindOneByID")
	})
	s.Run("internal-server-error", func() {
		s.SetupTest()
		s.servicesQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, serviceUrl, method).Return(services, nil)
		s.usersServicesQuery.On("FindOneByServiceIDAndUserId", mock.Anything, services.ID, userId).Return(&modelsUsersServices.UsersServices{}, nil)
		s.resourcesMapping.On("FindByServiceId", mock.Anything, services.ID).Return(nil, driver.ErrBadConn)
		res := s.usecases.GetDbsConnection(context.TODO(), userId, method, serviceUrl, "")
		s.NotNil(res.Error)
		s.resourcesMapping.AssertNotCalled(s.T(), "FindByServiceId")
		s.dbsPostgreQuery.AssertNotCalled(s.T(), "FindOneByID")
	})
	s.Run("internal-server-error", func() {
		s.SetupTest()
		s.servicesQuery.On("FindOneByServiceUrlAndMethodAndQueryIsNull", mock.Anything, serviceUrl, method).Return(services, nil)
		s.usersServicesQuery.On("FindOneByServiceIDAndUserId", mock.Anything, services.ID, userId).Return(&modelsUsersServices.UsersServices{}, nil)
		s.resourcesMapping.On("FindByServiceId", mock.Anything, services.ID).Return(modelsRM.ResourcesMappingList{modelsRM.ResourcesMapping{ID: "1"}}, nil)
		s.dbsPostgreQuery.On("FindOneByID", mock.Anything, services.DbID).Return(nil, driver.ErrBadConn)
		res := s.usecases.GetDbsConnection(context.TODO(), userId, method, serviceUrl, "")
		s.NotNil(res.Error)
		s.resourcesMapping.AssertNotCalled(s.T(), "FindByServiceId")
		s.dbsPostgreQuery.AssertNotCalled(s.T(), "FindOneByID")
	})
}
