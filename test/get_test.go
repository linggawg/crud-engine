package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"crud-engine/handler"
	conn "crud-engine/pkg/database"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	getJSON := `{"success":true,"data":{"content":[{"created_by":"system","created_date":"2021-10-15T21:09:06Z","id":"1","is_deleted":"0","name":"Aceh"}],"maxPage":null,"page":null,"size":null,"totalElements":1},"message":"List table province","code":200}
`

	sqlx, err := conn.InitSqlx(GlobalEnv.SQLXDatabase)
	if err != nil {
		panic(err)
	}
	log.Println("Database successfully initialized")

	q := make(url.Values)
	q.Set("query", "id = 1")
	q.Set("sortBy", "id")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("table")
	c.SetParamValues("province")

	if assert.NoError(t, handler.New(sqlx).Get(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, getJSON, rec.Body.String())
		log.Println("GET handler test success")
	}
}
