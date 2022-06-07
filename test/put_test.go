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

func TestPut(t *testing.T) {
	responsePutJSON := `{"success":true,"data":{"created_by":1,"created_on":"2020-04-30T17:35:27.638Z","is_deleted":false,"last_modified_by":1,"last_modified_on":"2020-04-30T17:35:27.638Z","province_name":"Kalimantan Timur"},"message":"successfully update province with province_id 11","code":200}`
	putJSON := `{
		"created_by": 1,
		"created_on": "2020-04-30T17:35:27.638Z",
		"is_deleted": false,
		"last_modified_by": 1,
		"last_modified_on": "2020-04-30T17:35:27.638Z",
		"province_name": "Kalimantan Timur"
	}`

	sqlx, err := conn.InitSqlx(TestEnv.SQLXDatabase)
	if err != nil {
		panic(err)
	}
	log.Println("Database successfully initialized")

	q := make(url.Values)
	q.Set("field_id", "province_id")

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/?"+q.Encode(), strings.NewReader(putJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	req.Header.Set(echo.HeaderAuthorization, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NTQ1OTQzMDgsImlhdCI6MTY1NDU5MDcwOCwic3ViIjoiMmI0YTg3MDYtY2Y0Ni00NTEzLWI0YmUtZTMwOWJkM2QyNjY1In0.ERHAfyyLl5PRjyDWYUgdj6ySrldKkrcfindWc6bX7xo")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("table", "value")
	c.SetParamValues("province", "11")

	if assert.NoError(t, handler.New(sqlx).Put(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, responsePutJSON, strings.TrimSpace(rec.Body.String()))
		log.Println("Put handler test success")
	}
}
