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

func TestDelete(t *testing.T) {
	responseDeleteJSON := `{"success":true,"data":0,"message":"successfully delete province with province_id 12","code":200}`
	
	sqlx, err := conn.InitSqlx(TestEnv.SQLXDatabase)
	if err != nil {
		panic(err)
	}
	log.Println("Database successfully initialized")
	
	q := make(url.Values)
	q.Set("field_id", "province_id")

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	req.Header.Set(echo.HeaderAuthorization, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NTQ1OTQzMDgsImlhdCI6MTY1NDU5MDcwOCwic3ViIjoiMmI0YTg3MDYtY2Y0Ni00NTEzLWI0YmUtZTMwOWJkM2QyNjY1In0.ERHAfyyLl5PRjyDWYUgdj6ySrldKkrcfindWc6bX7xo")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("table", "value")
	c.SetParamValues("province", "12")

	if assert.NoError(t, handler.New(sqlx).Delete(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, responseDeleteJSON, strings.TrimSpace(rec.Body.String()))
		log.Println("Delete handler test success")
	}
}