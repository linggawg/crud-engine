package handlers_test

import (
	"encoding/json"
	"engine/bin/modules/users/handlers"
	models "engine/bin/modules/users/models/domain"
	"engine/bin/modules/users/models/mocks"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestHttpSqlx(t *testing.T) {
	suite.Run(t, new(HttpSqlxTest))
}

type HttpSqlxTest struct {
	suite.Suite
	usecases *mocks.CommandUsecase
	handlers *handlers.HttpSqlx
}

func (s *HttpSqlxTest) SetupTest() {
	s.usecases = new(mocks.CommandUsecase)
	s.handlers = &handlers.HttpSqlx{
		CommandUsecase: s.usecases,
	}
}

func (s *HttpSqlxTest) TestLogin() {
	params := models.ReqLogin{
		Username: "admin",
		Password: "password",
		Duration: 3600,
	}
	paramsString, err := json.Marshal(params)
	assert.NoError(s.T(), err)
	e := echo.New()
	s.handlers.Mount(e.Group("/engine"))
	s.Run("success", func() {
		req := httptest.NewRequest(http.MethodPost, "/v1/login", strings.NewReader(string(paramsString)))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		s.usecases.On("Login", mock.AnythingOfType("*context.emptyCtx"), params).Return(utils.Result{Data: "token"})
		assert.NoError(s.T(), s.handlers.Login(c))
		assert.Equal(s.T(), http.StatusOK, rec.Code)
		s.usecases.AssertCalled(s.T(), "Login", mock.AnythingOfType("*context.emptyCtx"), params)
	})
	s.Run("error-login", func() {
		s.SetupTest()
		s.usecases.On("Login", mock.AnythingOfType("*context.emptyCtx"), params).Return(utils.Result{Error: httpError.NewBadRequest()})
		req := httptest.NewRequest(http.MethodPost, "/v1/login", strings.NewReader(string(paramsString)))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		assert.NoError(s.T(), s.handlers.Login(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.usecases.AssertCalled(s.T(), "Login", mock.AnythingOfType("*context.emptyCtx"), params)
	})
	s.Run("error-validator", func() {
		s.SetupTest()
		req := httptest.NewRequest(http.MethodPost, "/engine/v1/login", strings.NewReader(`{"username":"","password":"users-password","duration":3600}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		assert.NoError(s.T(), s.handlers.Login(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.usecases.AssertNotCalled(s.T(), "Login")
	})
	s.Run("error-decode", func() {
		s.SetupTest()
		req := httptest.NewRequest(http.MethodPost, "/engine/v1/login", strings.NewReader(`{"username":"","password":"users-password","duration":3600,}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		assert.NoError(s.T(), s.handlers.Login(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.usecases.AssertNotCalled(s.T(), "Login")
	})
}

func (s *HttpSqlxTest) TestRegisterUser() {
	params := models.ReqUser{
		Username: "usertest",
		Password: "usertest",
		RoleID:   "1",
	}
	paramsString, err := json.Marshal(params)
	assert.NoError(s.T(), err)
	e := echo.New()
	s.handlers.Mount(e.Group("/engine"))
	s.Run("success", func() {
		req := httptest.NewRequest(http.MethodPost, "/v1/register", strings.NewReader(string(paramsString)))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		s.usecases.On("RegisterUser", mock.AnythingOfType("*context.emptyCtx"), params).Return(utils.Result{Data: "token"})
		assert.NoError(s.T(), s.handlers.RegisterUser(c))
		assert.Equal(s.T(), http.StatusCreated, rec.Code)
		s.usecases.AssertCalled(s.T(), "RegisterUser", mock.AnythingOfType("*context.emptyCtx"), params)
	})
	s.Run("error-register", func() {
		s.SetupTest()
		s.usecases.On("RegisterUser", mock.AnythingOfType("*context.emptyCtx"), params).Return(utils.Result{Error: httpError.NewBadRequest()})
		req := httptest.NewRequest(http.MethodPost, "/v1/register", strings.NewReader(string(paramsString)))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		assert.NoError(s.T(), s.handlers.RegisterUser(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.usecases.AssertCalled(s.T(), "RegisterUser", mock.AnythingOfType("*context.emptyCtx"), params)
	})
	s.Run("error-validator", func() {
		s.SetupTest()
		req := httptest.NewRequest(http.MethodPost, "/v1/register", strings.NewReader(`{"username":"","password":"users-password","role_id":"1"}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		assert.NoError(s.T(), s.handlers.RegisterUser(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.usecases.AssertNotCalled(s.T(), "RegisterUser")
	})
	s.Run("error-decode", func() {
		s.SetupTest()
		req := httptest.NewRequest(http.MethodPost, "/v1/register", strings.NewReader(`{"username":"","password":"users-password","role_id":"1",}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		assert.NoError(s.T(), s.handlers.RegisterUser(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.usecases.AssertNotCalled(s.T(), "RegisterUser")
	})
}
