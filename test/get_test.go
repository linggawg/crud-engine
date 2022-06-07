package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"crud-engine/handler"
	conn "crud-engine/pkg/database"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	getJSON := `{"success":true,"data":{"content":[{"created_by":1,"created_on":"2020-04-29T00:00:00Z","is_deleted":false,"last_modified_by":1,"last_modified_on":"2020-04-29T00:00:00Z","province_id":1,"province_name":"Jawa Barat"}],"maxPage":null,"page":null,"size":null,"totalElements":1},"message":"List table province","code":200}`
	
	sqlx, err := conn.InitSqlx(TestEnv.SQLXDatabase)
	if err != nil {
		panic(err)
	}
	log.Println("Database successfully initialized")

	q := make(url.Values)
	q.Set("query", "province_id = 1")
	q.Set("sortBy", "province_id")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NTQ1OTQzMDgsImlhdCI6MTY1NDU5MDcwOCwic3ViIjoiMmI0YTg3MDYtY2Y0Ni00NTEzLWI0YmUtZTMwOWJkM2QyNjY1In0.ERHAfyyLl5PRjyDWYUgdj6ySrldKkrcfindWc6bX7xo")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("table")
	c.SetParamValues("province")

	if assert.NoError(t, handler.New(sqlx).Get(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, getJSON, strings.TrimSpace(rec.Body.String()))
		log.Println("GET handler test success")
	}
}
