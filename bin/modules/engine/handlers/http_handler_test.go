package handlers_test

import (
	"encoding/json"
	"engine/bin/config"
	modelsDbs "engine/bin/modules/dbs/models/domain"
	mocks2 "engine/bin/modules/dbs/models/mocks"
	"engine/bin/modules/engine/handlers"
	models "engine/bin/modules/engine/models/domain"
	"engine/bin/modules/engine/models/mocks"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/token"
	"engine/bin/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestEngineHTTPHandler(t *testing.T) {
	suite.Run(t, new(EngineHTTPHandlerTest))
}

type EngineHTTPHandlerTest struct {
	suite.Suite
	engineQueryUsecase   *mocks.EngineQueryUsecase
	engineCommandUsecase *mocks.EngineCommandUsecase
	dbsQueryUsecase      *mocks2.DbsQueryUsecase
	handlers             *handlers.EngineHTTPHandler
}

func (s *EngineHTTPHandlerTest) SetupTest() {
	s.engineQueryUsecase = new(mocks.EngineQueryUsecase)
	s.engineCommandUsecase = new(mocks.EngineCommandUsecase)
	s.dbsQueryUsecase = new(mocks2.DbsQueryUsecase)
	s.handlers = &handlers.EngineHTTPHandler{
		EngineQueryUsecase:   s.engineQueryUsecase,
		EngineCommandUsecase: s.engineCommandUsecase,
		DbsQueryUsecase:      s.dbsQueryUsecase,
	}
}
func (s *EngineHTTPHandlerTest) TestGet() {
	dbs := modelsDbs.Dbs{ID: "09582c5e-e7fd-49df-bfa4-a9428ef3a1b4", Name: "db_config", Host: "localhost", Port: 5432, Username: "root", Password: nil, Dialect: "mocks"}
	engineConfig := models.EngineConfig{Dbs: dbs}
	page := 0
	size := 10
	payload := &models.GetList{Page: &page, Size: &size, Sort: "", IsDistinct: true, Filter: "", Columns: "", Key: "", Opts: token.Claim{UserID: "1"}}
	table := "users"
	var data []map[string]interface{}
	data = append(data, map[string]interface{}{"id": "1", "name": "data1"}, map[string]interface{}{"id": "2", "name": "data3"})

	e := echo.New()
	s.handlers.Mount(e.Group("/engine"))

	s.Run("success", func() {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/:table", strings.NewReader(""))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table")
		c.SetParamValues(table)
		q := req.URL.Query()
		q.Add("key", payload.Key)
		q.Add("isDistinct", strconv.FormatBool(payload.IsDistinct))
		q.Add("pageNo", strconv.Itoa(*payload.Page))
		q.Add("pageSize", strconv.Itoa(*payload.Size))
		req.URL.RawQuery = q.Encode()
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), payload.Opts.UserID, req.Method, table, payload.Key).Return(utils.Result{Data: engineConfig})
		s.engineQueryUsecase.On("Get", mock.AnythingOfType("*context.emptyCtx"), engineConfig, table, payload).Return(utils.Result{Data: data})
		assert.NoError(s.T(), s.handlers.Get(c))
		assert.Equal(s.T(), http.StatusOK, rec.Code)
	})
	s.Run("error-get-dbs", func() {
		s.SetupTest()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/:table", strings.NewReader(""))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table")
		c.SetParamValues(table)
		q := req.URL.Query()
		q.Add("key", payload.Key)
		q.Add("isDistinct", strconv.FormatBool(payload.IsDistinct))
		q.Add("pageNo", strconv.Itoa(*payload.Page))
		q.Add("pageSize", strconv.Itoa(*payload.Size))
		req.URL.RawQuery = q.Encode()
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), payload.Opts.UserID, req.Method, table, payload.Key).Return(utils.Result{Error: httpError.NewNotFound()})
		assert.NoError(s.T(), s.handlers.Get(c))
		assert.Equal(s.T(), http.StatusNotFound, rec.Code)
		s.engineCommandUsecase.AssertNotCalled(s.T(), "Get")
	})
	s.Run("error", func() {
		s.SetupTest()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/:table", strings.NewReader(""))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table")
		c.SetParamValues(table)
		q := req.URL.Query()
		q.Add("key", payload.Key)
		q.Add("isDistinct", strconv.FormatBool(payload.IsDistinct))
		q.Add("pageNo", strconv.Itoa(*payload.Page))
		q.Add("pageSize", strconv.Itoa(*payload.Size))
		req.URL.RawQuery = q.Encode()
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), req.Method, table, "").Return(utils.Result{Data: engineConfig})
		s.engineQueryUsecase.On("Get", mock.AnythingOfType("*context.emptyCtx"), engineConfig, table, payload).Return(utils.Result{Error: httpError.NewInternalServerError()})
		assert.NoError(s.T(), s.handlers.Get(c))
		assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	})
}

func (s *EngineHTTPHandlerTest) TestPost() {
	dbs := modelsDbs.Dbs{ID: "09582c5e-e7fd-49df-bfa4-a9428ef3a1b4", Name: "db_config", Host: "localhost", Port: 5432, Username: "root", Password: nil, Dialect: "mocks"}
	engineConfig := models.EngineConfig{Dbs: dbs}
	payload := &models.EngineRequest{Table: "users", Data: map[string]interface{}{"id": "1", "name": "data1"}, Opts: token.Claim{UserID: "1"}}
	payloadString, err := json.Marshal(payload.Data)
	assert.NoError(s.T(), err)
	table := "users"
	e := echo.New()
	s.handlers.Mount(e.Group("/engine"))
	s.Run("success", func() {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/:table", strings.NewReader(string(payloadString)))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table")
		c.SetParamValues(table)
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), req.Method, table, "").Return(utils.Result{Data: engineConfig})
		s.engineCommandUsecase.On("Insert", mock.AnythingOfType("*context.emptyCtx"), engineConfig, payload).Return(utils.Result{Data: payload})
		assert.NoError(s.T(), s.handlers.Post(c))
		assert.Equal(s.T(), http.StatusOK, rec.Code)
	})
	s.Run("error-decode", func() {
		s.SetupTest()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/:table", strings.NewReader(`{"id": "1", "name": "data1",}`))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table")
		c.SetParamValues(table)
		assert.NoError(s.T(), s.handlers.Post(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.engineCommandUsecase.AssertNotCalled(s.T(), "Insert")
	})
	s.Run("error-get-dbs", func() {
		s.SetupTest()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/:table", strings.NewReader(string(payloadString)))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table")
		c.SetParamValues(table)
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), req.Method, table, "").Return(utils.Result{Error: httpError.NewInternalServerError()})
		assert.NoError(s.T(), s.handlers.Post(c))
		assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
		s.engineCommandUsecase.AssertNotCalled(s.T(), "Insert")
	})
	s.Run("error", func() {
		s.SetupTest()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/:table", strings.NewReader(string(payloadString)))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table")
		c.SetParamValues(table)
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), req.Method, table, "").Return(utils.Result{Data: engineConfig})
		s.engineCommandUsecase.On("Insert", mock.AnythingOfType("*context.emptyCtx"), engineConfig, payload).Return(utils.Result{Error: httpError.NewConflict()})
		assert.NoError(s.T(), s.handlers.Post(c))
		assert.Equal(s.T(), http.StatusConflict, rec.Code)
	})
}
func (s *EngineHTTPHandlerTest) TestPut() {
	dbs := modelsDbs.Dbs{ID: "09582c5e-e7fd-49df-bfa4-a9428ef3a1b4", Name: "db_config", Host: "localhost", Port: 5432, Username: "root", Password: nil, Dialect: "mocks"}
	engineConfig := models.EngineConfig{Dbs: dbs}
	payload := &models.EngineRequest{Table: "users", FieldId: "id", Value: "1", Data: map[string]interface{}{"id": "1", "name": "data1"}, Opts: token.Claim{UserID: "1"}}
	payloadString, err := json.Marshal(payload.Data)
	assert.NoError(s.T(), err)
	table := "users"
	e := echo.New()
	s.handlers.Mount(e.Group("/engine"))
	config.GlobalEnv.StrictRestfulAPI = true
	s.Run("success", func() {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/v1/:table/:value", strings.NewReader(string(payloadString)))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table", "value")
		c.SetParamValues(table, payload.Value)
		q := req.URL.Query()
		q.Add("field_id", payload.FieldId)
		req.URL.RawQuery = q.Encode()
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), req.Method, table, "").Return(utils.Result{Data: engineConfig})
		s.engineCommandUsecase.On("Update", mock.AnythingOfType("*context.emptyCtx"), engineConfig, payload).Return(utils.Result{Data: payload})
		assert.NoError(s.T(), s.handlers.Put(c))
		assert.Equal(s.T(), http.StatusOK, rec.Code)
	})
	s.Run("error-decode", func() {
		s.SetupTest()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/v1/:table/:value", strings.NewReader(`{"id": "1", "name": "data1",}`))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table", "value")
		c.SetParamValues(table, payload.Value)
		q := req.URL.Query()
		q.Add("field_id", payload.FieldId)
		req.URL.RawQuery = q.Encode()
		assert.NoError(s.T(), s.handlers.Put(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.engineCommandUsecase.AssertNotCalled(s.T(), "Update")
	})
	s.Run("error-get-dbs", func() {
		s.SetupTest()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/v1/:table/:value", strings.NewReader(string(payloadString)))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table", "value")
		c.SetParamValues(table, payload.Value)
		q := req.URL.Query()
		q.Add("field_id", payload.FieldId)
		req.URL.RawQuery = q.Encode()
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), req.Method, table, "").Return(utils.Result{Error: httpError.NewInternalServerError()})
		assert.NoError(s.T(), s.handlers.Put(c))
		assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
		s.engineCommandUsecase.AssertNotCalled(s.T(), "Update")
	})
	s.Run("error", func() {
		s.SetupTest()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/v1/:table/:value", strings.NewReader(string(payloadString)))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table", "value")
		c.SetParamValues(table, payload.Value)
		q := req.URL.Query()
		q.Add("field_id", payload.FieldId)
		req.URL.RawQuery = q.Encode()
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), req.Method, table, "").Return(utils.Result{Data: engineConfig})
		config.GlobalEnv.StrictRestfulAPI = false
		s.engineCommandUsecase.On("Patch", mock.AnythingOfType("*context.emptyCtx"), engineConfig, payload).Return(utils.Result{Error: httpError.NewUnprocessableEntity()})
		assert.NoError(s.T(), s.handlers.Put(c))
		assert.Equal(s.T(), http.StatusUnprocessableEntity, rec.Code)
	})
}

func (s *EngineHTTPHandlerTest) TestPatch() {
	dbs := modelsDbs.Dbs{ID: "09582c5e-e7fd-49df-bfa4-a9428ef3a1b4", Name: "db_config", Host: "localhost", Port: 5432, Username: "root", Password: nil, Dialect: "mocks"}
	engineConfig := models.EngineConfig{Dbs: dbs}
	payload := &models.EngineRequest{Table: "users", FieldId: "id", Value: "1", Data: map[string]interface{}{"id": "1", "name": "data1"}, Opts: token.Claim{UserID: "1"}}
	payloadString, err := json.Marshal(payload.Data)
	assert.NoError(s.T(), err)
	table := "users"
	e := echo.New()
	s.handlers.Mount(e.Group("/engine"))
	s.Run("success", func() {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/v1/:table/:value", strings.NewReader(string(payloadString)))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table", "value")
		c.SetParamValues(table, payload.Value)
		q := req.URL.Query()
		q.Add("field_id", payload.FieldId)
		req.URL.RawQuery = q.Encode()
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), req.Method, table, "").Return(utils.Result{Data: engineConfig})
		s.engineCommandUsecase.On("Patch", mock.AnythingOfType("*context.emptyCtx"), engineConfig, payload).Return(utils.Result{Data: payload})
		assert.NoError(s.T(), s.handlers.Patch(c))
		assert.Equal(s.T(), http.StatusOK, rec.Code)
	})
	s.Run("error-decode", func() {
		s.SetupTest()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/v1/:table/:value", strings.NewReader(`{"id": "1", "name": "data1",}`))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table", "value")
		c.SetParamValues(table, payload.Value)
		q := req.URL.Query()
		q.Add("field_id", payload.FieldId)
		req.URL.RawQuery = q.Encode()
		assert.NoError(s.T(), s.handlers.Patch(c))
		assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
		s.engineCommandUsecase.AssertNotCalled(s.T(), "Patch")
	})
	s.Run("error-get-dbs", func() {
		s.SetupTest()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/v1/:table/:value", strings.NewReader(string(payloadString)))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table", "value")
		c.SetParamValues(table, payload.Value)
		q := req.URL.Query()
		q.Add("field_id", payload.FieldId)
		req.URL.RawQuery = q.Encode()
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), req.Method, table, "").Return(utils.Result{Error: httpError.NewUnauthorized()})
		assert.NoError(s.T(), s.handlers.Patch(c))
		assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
		s.engineCommandUsecase.AssertNotCalled(s.T(), "Patch")
	})
	s.Run("error", func() {
		s.SetupTest()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/v1/:table/:value", strings.NewReader(string(payloadString)))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table", "value")
		c.SetParamValues(table, payload.Value)
		q := req.URL.Query()
		q.Add("field_id", payload.FieldId)
		req.URL.RawQuery = q.Encode()
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), req.Method, table, "").Return(utils.Result{Data: engineConfig})
		s.engineCommandUsecase.On("Patch", mock.AnythingOfType("*context.emptyCtx"), engineConfig, payload).Return(utils.Result{Error: httpError.NewInternalServerError()})
		assert.NoError(s.T(), s.handlers.Patch(c))
		assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	})
}

func (s *EngineHTTPHandlerTest) TestDelete() {
	dbs := modelsDbs.Dbs{ID: "09582c5e-e7fd-49df-bfa4-a9428ef3a1b4", Name: "db_config", Host: "localhost", Port: 5432, Username: "root", Password: nil, Dialect: "mocks"}
	engineConfig := models.EngineConfig{Dbs: dbs}
	payload := &models.EngineRequest{Table: "users", FieldId: "id", Value: "1", Opts: token.Claim{UserID: "1"}}
	table := "users"
	e := echo.New()
	s.handlers.Mount(e.Group("/engine"))
	s.Run("success", func() {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/v1/:table/:value", strings.NewReader(""))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table", "value")
		c.SetParamValues(table, payload.Value)
		q := req.URL.Query()
		q.Add("field_id", payload.FieldId)
		req.URL.RawQuery = q.Encode()
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), req.Method, table, "").Return(utils.Result{Data: engineConfig})
		s.engineCommandUsecase.On("Delete", mock.AnythingOfType("*context.emptyCtx"), engineConfig, payload).Return(utils.Result{Data: payload})
		assert.NoError(s.T(), s.handlers.Delete(c))
		assert.Equal(s.T(), http.StatusOK, rec.Code)
	})
	s.Run("error-get-dbs", func() {
		s.SetupTest()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/v1/:table/:value", strings.NewReader(""))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table", "value")
		c.SetParamValues(table, payload.Value)
		q := req.URL.Query()
		q.Add("field_id", payload.FieldId)
		req.URL.RawQuery = q.Encode()
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), req.Method, table, "").Return(utils.Result{Error: httpError.CommonError{Code: http.StatusConflict, Message: "Conflict"}})
		assert.NoError(s.T(), s.handlers.Delete(c))
		assert.Equal(s.T(), http.StatusConflict, rec.Code)
		s.engineCommandUsecase.AssertNotCalled(s.T(), "Delete")
	})
	s.Run("error", func() {
		s.SetupTest()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/v1/:table/:value", strings.NewReader(""))
		c := e.NewContext(req, rec)
		c.Set("opts", payload.Opts)
		c.SetParamNames("table", "value")
		c.SetParamValues(table, payload.Value)
		q := req.URL.Query()
		q.Add("field_id", payload.FieldId)
		req.URL.RawQuery = q.Encode()
		s.dbsQueryUsecase.On("GetDbsConnection", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), req.Method, table, "").Return(utils.Result{Data: engineConfig})
		s.engineCommandUsecase.On("Delete", mock.AnythingOfType("*context.emptyCtx"), engineConfig, payload).Return(utils.Result{Error: httpError.NewInternalServerError()})
		assert.NoError(s.T(), s.handlers.Delete(c))
		assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	})
}
