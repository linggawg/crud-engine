package queries_test

import (
	"context"
	"database/sql"
	models "engine/bin/modules/services/models/domain"
	"engine/bin/modules/services/repositories/queries"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

var (
	serviceDefinition = "SELECT id FROM user"
	serviceUrl        = "user"
	mockServicesList  = models.ServicesList{
		models.Services{
			ID:                "20bd16be-7f64-4281-9b58-758ecc7098b9",
			DbID:              "33d2a2b4-463f-4a27-aaac-d90fbfe6ebd2",
			Method:            "GET",
			ServiceUrl:        &serviceUrl,
			ServiceDefinition: &serviceDefinition,
			IsQuery:           false,
			CreatedAt:         null.TimeFrom(time.Now()),
			CreatedBy:         "93147c1d-197e-4556-9d24-1992196aaa03",
			ModifiedAt:        null.TimeFrom(time.Now()),
			ModifiedBy:        "93147c1d-197e-4556-9d24-1992196aaa03",
		},
		models.Services{
			ID:                "18a53750-d3e2-4e5a-9d0c-e88d368509f7",
			DbID:              "33d2a2b4-463f-4a27-aaac-d90fbfe6ebd2",
			Method:            "GET",
			ServiceUrl:        &serviceUrl,
			ServiceDefinition: &serviceDefinition,
			IsQuery:           false,
			CreatedAt:         null.TimeFrom(time.Now()),
			CreatedBy:         "93147c1d-197e-4556-9d24-1992196aaa03",
			ModifiedAt:        null.TimeFrom(time.Now()),
			ModifiedBy:        "93147c1d-197e-4556-9d24-1992196aaa03",
		},
	}
)

func TestGetByServiceUrlAndMethod(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databases connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "db_id", "method", "service_url", "service_definition", "is_query", "created_at", "created_by", "modified_at", "modified_by"})
	for _, i := range mockServicesList {
		rows.AddRow(i.ID, i.DbID, i.Method, i.ServiceUrl, i.ServiceDefinition, i.IsQuery, i.CreatedAt, i.CreatedBy, i.ModifiedAt, i.ModifiedBy)
	}

	query := "SELECT id, db_id, service_url, method, service_definition, is_query FROM services WHERE service_url = \\$1 AND method = \\$2"
	mock.ExpectQuery(query).WithArgs(mockServicesList[0].ServiceUrl, mockServicesList[0].Method).WillReturnRows(rows)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	a := queries.NewServicesQuery(sqlxDB)

	mService, err := a.GetByServiceUrlAndMethod(context.TODO(), *mockServicesList[0].ServiceUrl, mockServicesList[0].Method)
	assert.NoError(t, err)
	assert.NotNil(t, mService)
}

func TestGetByServiceUrlAndMethodNoRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databases connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "db_id", "method", "service_url", "service_definition", "is_query", "created_at", "created_by", "modified_at", "modified_by"})
	query := "SELECT id, db_id, service_url, method, service_definition, is_query FROM services WHERE service_url = \\$1 AND method = \\$2"
	mock.ExpectQuery(query).WithArgs(mockServicesList[0].ServiceUrl, mockServicesList[0].Method).WillReturnRows(rows)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	a := queries.NewServicesQuery(sqlxDB)

	mService, err := a.GetByServiceUrlAndMethod(context.TODO(), *mockServicesList[0].ServiceUrl, mockServicesList[0].Method)
	assert.Empty(t, mService)
	assert.Error(t, err)
	assert.Equal(t, err, sql.ErrNoRows)
}

func TestGetByServiceDefinitionAndMethod(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databases connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "db_id", "method", "service_url", "service_definition", "is_query", "created_at", "created_by", "modified_at", "modified_by"})
	for _, i := range mockServicesList {
		rows.AddRow(i.ID, i.DbID, i.Method, i.ServiceUrl, i.ServiceDefinition, i.IsQuery, i.CreatedAt, i.CreatedBy, i.ModifiedAt, i.ModifiedBy)
	}

	query := "SELECT id, db_id, service_url, method, service_definition, is_query FROM services WHERE is_query = true AND service_definition = \\$1 AND method = \\$2"
	mock.ExpectQuery(query).WithArgs(mockServicesList[0].ServiceDefinition, mockServicesList[0].Method).WillReturnRows(rows)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	a := queries.NewServicesQuery(sqlxDB)

	mService, err := a.GetByServiceDefinitionAndMethod(context.TODO(), *mockServicesList[0].ServiceDefinition, mockServicesList[0].Method)
	assert.NoError(t, err)
	assert.NotNil(t, mService)
}

func TestGetByServiceDefinitionAndMethodNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databases connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "db_id", "method", "service_url", "service_definition", "is_query", "created_at", "created_by", "modified_at", "modified_by"})
	query := "SELECT id, db_id, service_url, method, service_definition, is_query FROM services WHERE is_query = true AND service_definition = \\$1 AND method = \\$2"
	mock.ExpectQuery(query).WithArgs(mockServicesList[0].ServiceDefinition, mockServicesList[0].Method).WillReturnRows(rows)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	a := queries.NewServicesQuery(sqlxDB)

	mService, err := a.GetByServiceDefinitionAndMethod(context.TODO(), *mockServicesList[0].ServiceDefinition, mockServicesList[0].Method)
	assert.Empty(t, mService)
	assert.Error(t, err)
	assert.Equal(t, err, sql.ErrNoRows)
}
