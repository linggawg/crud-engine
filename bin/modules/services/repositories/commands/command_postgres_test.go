package commands_test

import (
	"context"
	models "engine/bin/modules/services/models/domain"
	"engine/bin/modules/services/repositories/commands"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func NewMock() (*sqlx.DB, sqlmock.Sqlmock) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db := sqlx.NewDb(mockDb, "sqlmock")
	return db, mock
}

var random = uuid.New().String()
var services = &models.Services{
	ID:         random,
	DbID:       random,
	QueryID:    &random,
	Method:     random,
	ServiceUrl: &random,
	CreatedAt:  null.TimeFrom(time.Now()),
	CreatedBy:  &random,
	ModifiedAt: null.TimeFrom(time.Now()),
	ModifiedBy: &random,
}

func TestInsertOne(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock := NewMock()

		query := "INSERT INTO services ( id, db_id, query_id, method, service_url, created_at, created_by, modified_at, modified_by ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? );"

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(services.ID, services.DbID, services.QueryID, services.Method, services.ServiceUrl, services.CreatedAt, services.CreatedBy, services.ModifiedAt, services.ModifiedBy).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := commands.NewServicesCommand(db).InsertOne(context.TODO(), services)
		assert.NoError(t, err)
	})

	t.Run("Error Begin", func(t *testing.T) {
		db, mock := NewMock()

		query := "INSERT INTO services ( id, db_id, query_id, method, service_url, created_at, created_by, modified_at, modified_by ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? );"

		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(services.ID, services.DbID, services.QueryID, services.Method, services.ServiceUrl, services.CreatedAt, services.CreatedBy, services.ModifiedAt, services.ModifiedBy).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := commands.NewServicesCommand(db).InsertOne(context.TODO(), services)
		assert.Error(t, err)
	})
	t.Run("Error Exec", func(t *testing.T) {
		db, mock := NewMock()

		query := "INSERT INTO services ( id, db_id, query_id, method, service_url, created_at, created_by, modified_at, modified_by ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? );"

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(services.ID, nil, services.QueryID, services.Method, services.ServiceUrl, services.CreatedAt, services.CreatedBy, services.ModifiedAt, services.ModifiedBy).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := commands.NewServicesCommand(db).InsertOne(context.TODO(), services)
		assert.Error(t, err)
	})
	t.Run("Error Commit", func(t *testing.T) {
		db, mock := NewMock()

		query := "INSERT INTO services ( id, db_id, query_id, method, service_url, created_at, created_by, modified_at, modified_by ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? );"

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(services.ID, services.DbID, services.QueryID, services.Method, services.ServiceUrl, services.CreatedAt, services.CreatedBy, services.ModifiedAt, services.ModifiedBy).WillReturnResult(sqlmock.NewResult(0, 0))

		err := commands.NewServicesCommand(db).InsertOne(context.TODO(), services)
		assert.Error(t, err)
	})
	t.Run("DB no connection", func(t *testing.T) {
		db, _ := NewMock()

		db.Close()
		err := commands.NewServicesCommand(db).InsertOne(context.TODO(), services)
		assert.Error(t, err)
		assert.EqualError(t, err, "error establishing a database connection")
	})
}

func TestDeleteByServiceUrl(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock := NewMock()

		query := "DELETE FROM services WHERE service_url = $1;"

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(services.ServiceUrl).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := commands.NewServicesCommand(db).DeleteByServiceUrl(context.TODO(), *services.ServiceUrl)
		assert.NoError(t, err)
	})

	t.Run("Error Begin", func(t *testing.T) {
		db, mock := NewMock()

		query := "DELETE FROM services WHERE service_url = $1;"

		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(services.ServiceUrl).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := commands.NewServicesCommand(db).DeleteByServiceUrl(context.TODO(), *services.ServiceUrl)
		assert.Error(t, err)
	})
	t.Run("Error Exec", func(t *testing.T) {
		db, mock := NewMock()

		query := "DELETE FROM services WHERE service_url = $1;"

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(nil).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := commands.NewServicesCommand(db).DeleteByServiceUrl(context.TODO(), *services.ServiceUrl)
		assert.Error(t, err)
	})
	t.Run("Error Commit", func(t *testing.T) {
		db, mock := NewMock()

		query := "DELETE FROM services WHERE service_url = $1;"

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(services.ServiceUrl).WillReturnResult(sqlmock.NewResult(0, 0))

		err := commands.NewServicesCommand(db).DeleteByServiceUrl(context.TODO(), *services.ServiceUrl)
		assert.Error(t, err)
	})
	t.Run("DB no connection", func(t *testing.T) {
		db, _ := NewMock()

		db.Close()
		err := commands.NewServicesCommand(db).DeleteByServiceUrl(context.TODO(), *services.ServiceUrl)
		assert.Error(t, err)
		assert.EqualError(t, err, "error establishing a database connection")
	})
}
