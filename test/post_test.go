package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"crud-engine/handler"
	conn "crud-engine/pkg/database"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	responsePostJSON := `{"success":true,"data":{"created_by":1,"created_on":"2020-04-30T17:35:27.638Z","is_deleted":false,"last_modified_by":1,"last_modified_on":"2020-04-30T17:35:27.638Z","province_name":"Kalimantan Timur"},"message":"successfully insert province","code":201}`
	postJSON := `{
		"is_deleted": false,
		"province_name": "Kalimantan Timur",
		"created_by": 1,
		"created_on": "2020-04-30T17:35:27.638Z",
		"last_modified_by": 1,
		"last_modified_on": "2020-04-30T17:35:27.638Z"
	}`
	
	sqlx, err := conn.InitSqlx(TestEnv.SQLXDatabase)
	if err != nil {
		panic(err)
	}
	log.Println("Database successfully initialized")
	
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(postJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	req.Header.Set(echo.HeaderAuthorization, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NTQ1OTQzMDgsImlhdCI6MTY1NDU5MDcwOCwic3ViIjoiMmI0YTg3MDYtY2Y0Ni00NTEzLWI0YmUtZTMwOWJkM2QyNjY1In0.ERHAfyyLl5PRjyDWYUgdj6ySrldKkrcfindWc6bX7xo")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("table")
	c.SetParamValues("province")
	
	if assert.NoError(t, handler.New(sqlx).Post(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, responsePostJSON, strings.TrimSpace(rec.Body.String()))
		log.Println("POST handler test success")
	}
}