package handlers_test

import (
	"encoding/json"
	"engine/bin/modules/users-services/handlers"
	models "engine/bin/modules/users-services/models/domain"
	"engine/bin/modules/users-services/models/mocks"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/token"
	"engine/bin/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHttpSqlx(t *testing.T) {
	suite.Run(t, new(HttpSqlxTest))
}

type HttpSqlxTest struct {
	suite.Suite
	usecases *mocks.ServicesCommandUsecase
	handlers *handlers.HttpSqlx
}

func (s *HttpSqlxTest) SetupTest() {
	s.usecases = new(mocks.ServicesCommandUsecase)
	s.handlers = &handlers.HttpSqlx{
		CommandUsecase: s.usecases,
	}
}

func (s *HttpSqlxTest) TestGenerateDefaultUsersServices() {
	payload := models.UsersServicesRequest{
		ServiceUrl: "users",
		UserID:     "a14e2f3a-8884-4f97-9216-6620081f130e",
		Opts:       token.Claim{},
	}
	payloadString, err := json.Marshal(payload)
	assert.NoError(s.T(), err)
	e := echo.New()
	s.handlers.Mount(e.Group("/engine"))
	s.Run("success", func() {
		req := httptest.NewRequest(http.MethodPost, "/v1/users-services/default", strings.NewReader(string(payloadString)))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		s.usecases.On("InsertAllByServices", mock.AnythingOfType("*context.emptyCtx"), payload).Return(utils.Result{Data: "Generate Default Users Services Method"})
		assert.NoError(s.T(), s.handlers.GenerateDefaultUsersServices(c))
		assert.Equal(s.T(), http.StatusCreated, rec.Code)
		s.usecases.AssertCalled(s.T(), "InsertAllByServices", mock.AnythingOfType("*context.emptyCtx"), payload)
	})
	s.Run("error", func() {
		s.SetupTest()
		req := httptest.NewRequest(http.MethodPost, "/v1/users-services/default", strings.NewReader(string(payloadString)))
		rec := httptest.NewRecorder()
		s.usecases.On("InsertAllByServices", mock.AnythingOfType("*context.emptyCtx"), payload).Return(utils.Result{Error: httpError.NewBadRequest()})
		c := e.NewContext(req, rec)
		assert.NoError(s.T(), s.handlers.GenerateDefaultUsersServices(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.usecases.AssertCalled(s.T(), "InsertAllByServices", mock.AnythingOfType("*context.emptyCtx"), payload)
	})
	s.Run("error-validator", func() {
		s.SetupTest()
		req := httptest.NewRequest(http.MethodPost, "/v1/users-services/default", strings.NewReader(`{"service_url":"","opts":null}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		assert.NoError(s.T(), s.handlers.GenerateDefaultUsersServices(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.usecases.AssertNotCalled(s.T(), "InsertAllByServices")
	})
	s.Run("error-decode", func() {
		s.SetupTest()
		req := httptest.NewRequest(http.MethodPost, "/v1/users-services/default", strings.NewReader(`{"service_url":"","opts":""}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		assert.NoError(s.T(), s.handlers.GenerateDefaultUsersServices(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.usecases.AssertNotCalled(s.T(), "InsertAllByServices")
	})
}

func (s *HttpSqlxTest) TestDeleteByUsersIdAndServiceUrl() {
	payload := models.UsersServicesRequest{
		ServiceUrl: "users",
		UserID:     "a14e2f3a-8884-4f97-9216-6620081f130e",
		Opts: token.Claim{
			UserID:        "",
			RoleName:      "Admin",
			Authorization: "",
		},
	}
	payloadString, err := json.Marshal(payload)
	assert.NoError(s.T(), err)
	e := echo.New()
	s.handlers.Mount(e.Group("/engine"))
	s.Run("success", func() {
		req := httptest.NewRequest(http.MethodDelete, "/v1/users-services/default", strings.NewReader(string(payloadString)))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		s.usecases.On("DeleteByUsersIdAndServiceUrl", mock.AnythingOfType("*context.emptyCtx"), payload).Return(utils.Result{Data: "Delete Default Users Services Method"})
		assert.NoError(s.T(), s.handlers.DeleteByUsersIdAndServiceUrl(c))
		assert.Equal(s.T(), http.StatusOK, rec.Code)
		s.usecases.AssertCalled(s.T(), "DeleteByUsersIdAndServiceUrl", mock.AnythingOfType("*context.emptyCtx"), payload)
	})
	s.Run("error", func() {
		s.SetupTest()
		req := httptest.NewRequest(http.MethodDelete, "/v1/users-services/default", strings.NewReader(string(payloadString)))
		rec := httptest.NewRecorder()
		s.usecases.On("DeleteByUsersIdAndServiceUrl", mock.AnythingOfType("*context.emptyCtx"), payload).Return(utils.Result{Error: httpError.NewBadRequest()})
		c := e.NewContext(req, rec)
		assert.NoError(s.T(), s.handlers.DeleteByUsersIdAndServiceUrl(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.usecases.AssertCalled(s.T(), "DeleteByUsersIdAndServiceUrl", mock.AnythingOfType("*context.emptyCtx"), payload)
	})
	s.Run("error-validator", func() {
		s.SetupTest()
		req := httptest.NewRequest(http.MethodDelete, "/v1/users-services/default", strings.NewReader(`{"service_url":"","opts":null}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		assert.NoError(s.T(), s.handlers.DeleteByUsersIdAndServiceUrl(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.usecases.AssertNotCalled(s.T(), "DeleteByUsersIdAndServiceUrl")
	})
	s.Run("error-decode", func() {
		s.SetupTest()
		req := httptest.NewRequest(http.MethodDelete, "/v1/users-services/default", strings.NewReader(`{"service_url":"","opts":""}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		assert.NoError(s.T(), s.handlers.DeleteByUsersIdAndServiceUrl(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.usecases.AssertNotCalled(s.T(), "DeleteByUsersIdAndServiceUrl")
	})
}
