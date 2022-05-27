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

var (
)

func TestPost(t *testing.T) {
	responsePostJSON := `{"success":true,"data":{"id":"99","is_deleted":"0","name":"TEST POST PROVINCE"},"message":"successfully insert province","code":201}
`
	postJSON := `{"id": "99","is_deleted": "0","name": "TEST POST PROVINCE"}`
	
	sqlx, err := conn.InitSqlx(GlobalEnv.SQLXDatabase)
	if err != nil {
		panic(err)
	}
	log.Println("Database successfully initialized")
	
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(postJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("table")
	c.SetParamValues("province")
	
	if assert.NoError(t, handler.New(sqlx).Post(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, responsePostJSON, rec.Body.String())
		log.Println("POST handler test success")
	}
}