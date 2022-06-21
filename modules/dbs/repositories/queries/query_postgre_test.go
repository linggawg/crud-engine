package queries_test

import (
	"context"
	"crud-engine/modules/dbs/models/domain"
	"crud-engine/modules/dbs/repositories/queries"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"testing"
	"time"
)

var userId = "48c8fe7a-49fa-48c5-8e80-d080146e1907"
var mockDbsList = models.DbsList{
	models.Dbs{
		ID:         "f6cf97f7-c803-4cfa-b1b1-0e4f57127d43",
		AppID:      "766c2de8-cbbe-4962-950d-b7ee3455d05f",
		Name:       "dbprovinsi",
		Host:       "localhost",
		Port:       3306,
		Username:   "root",
		Password:   nil,
		Dialect:    "mysql",
		CreatedAt:  null.TimeFrom(time.Now()),
		CreatedBy:  &userId,
		ModifiedAt: null.TimeFrom(time.Now()),
		ModifiedBy: &userId,
	},
	models.Dbs{
		ID:         "f6a987dd-63fa-4d5e-acc9-ddec7dc836be",
		AppID:      "766c2de8-cbbe-4962-950d-b7ee3455d05f",
		Name:       "crud_go",
		Host:       "localhost",
		Port:       5432,
		Username:   "postgres",
		Password:   nil,
		Dialect:    "postgres",
		CreatedAt:  null.TimeFrom(time.Now()),
		CreatedBy:  &userId,
		ModifiedAt: null.TimeFrom(time.Now()),
		ModifiedBy: &userId,
	},
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "app_id", "name", "host", "port", "username", "password", "dialect", "created_at", "created_by", "modified_at", "modified_by"})
	for _, d := range mockDbsList {
		rows.AddRow(d.ID, d.AppID, d.Name, d.Host, d.Port, d.Username, d.Password, d.Dialect, d.CreatedAt, d.CreatedBy, d.ModifiedAt, d.ModifiedBy)
	}

	query := "SELECT id, app_id, name, host, port, username, password, dialect FROM dbs WHERE id = \\$1"
	mock.ExpectQuery(query).WithArgs(mockDbsList[0].ID).WillReturnRows(rows)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	a := queries.NewDbsQuery(sqlxDB)

	mUserServices, err := a.GetByID(context.TODO(), mockDbsList[0].ID)
	assert.NoError(t, err)
	assert.NotNil(t, mUserServices)
}

func TestGetByIDNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "app_id", "name", "host", "port", "username", "password", "dialect", "created_at", "created_by", "modified_at", "modified_by"})
	query := "SELECT id, app_id, name, host, port, username, password, dialect FROM dbs WHERE id = \\$1"
	mock.ExpectQuery(query).WithArgs(mockDbsList[0].ID).WillReturnRows(rows)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	a := queries.NewDbsQuery(sqlxDB)

	mDbs, err := a.GetByID(context.TODO(), mockDbsList[0].ID)
	assert.Empty(t, mDbs)
	assert.Error(t, err)
	assert.Equal(t, err, sql.ErrNoRows)
}
