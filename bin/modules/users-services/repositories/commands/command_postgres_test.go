package commands_test

import (
	"context"
	models "engine/bin/modules/users-services/models/domain"
	"engine/bin/modules/users-services/repositories/commands"
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
var usersServices = &models.UsersServices{
	ID:         random,
	UserID:     random,
	ServiceID:  random,
	CreatedAt:  null.TimeFrom(time.Now()),
	CreatedBy:  &random,
	ModifiedAt: null.TimeFrom(time.Now()),
	ModifiedBy: &random,
}

func TestDeleteByUsersIdAndServiceUrl(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock := NewMock()

		query := regexp.QuoteMeta("DELETE FROM users_services WHERE user_id = $1 AND service_id IN (SELECT id FROM services WHERE service_url = $2);")

		mock.ExpectBegin()
		mock.ExpectExec(query).WithArgs(usersServices.ID, "mock service url").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := commands.NewUsersServicesCommand(db).DeleteByUsersIdAndServiceUrl(context.TODO(), usersServices.ID, "mock service url")
		assert.NoError(t, err)
	})
	t.Run("Error Begin", func(t *testing.T) {
		db, mock := NewMock()

		query := regexp.QuoteMeta("DELETE FROM users_services WHERE user_id = $1 AND service_id IN (SELECT id FROM services WHERE service_url = $2);")

		mock.ExpectExec(query).WithArgs(usersServices.ID, "mock service url").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := commands.NewUsersServicesCommand(db).DeleteByUsersIdAndServiceUrl(context.TODO(), usersServices.ID, "mock service url")
		assert.Error(t, err)
	})
	t.Run("Error Exec", func(t *testing.T) {
		db, mock := NewMock()

		query := regexp.QuoteMeta("DELETE FROM users_services WHERE user_id = $1 AND service_id IN (SELECT id FROM services WHERE service_url = $2);")

		mock.ExpectBegin()
		mock.ExpectExec(query).WithArgs(nil, "mock service url").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := commands.NewUsersServicesCommand(db).DeleteByUsersIdAndServiceUrl(context.TODO(), usersServices.ID, "mock service url")
		assert.Error(t, err)
	})
	t.Run("Error Commit", func(t *testing.T) {
		db, mock := NewMock()

		query := regexp.QuoteMeta("DELETE FROM users_services WHERE user_id = $1 AND service_id IN (SELECT id FROM services WHERE service_url = $2);")

		mock.ExpectBegin()
		mock.ExpectExec(query).WithArgs(usersServices.ID, "mock service url").WillReturnResult(sqlmock.NewResult(0, 0))

		err := commands.NewUsersServicesCommand(db).DeleteByUsersIdAndServiceUrl(context.TODO(), usersServices.ID, "mock service url")
		assert.Error(t, err)
	})
	t.Run("DB no connection", func(t *testing.T) {
		db, _ := NewMock()

		db.Close()
		err := commands.NewUsersServicesCommand(db).DeleteByUsersIdAndServiceUrl(context.TODO(), usersServices.ID, "mock service url")
		assert.Error(t, err)
		assert.EqualError(t, err, "error establishing a database connection")
	})
}

func TestInsertOne(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock := NewMock()

		query := regexp.QuoteMeta("INSERT INTO users_services ( id, user_id, service_id, created_at, created_by, modified_at, modified_by ) VALUES ( ?, ?, ?, ?, ?, ?, ? );")

		mock.ExpectBegin()
		mock.ExpectExec(query).WithArgs(usersServices.ID, usersServices.UserID, usersServices.ServiceID, usersServices.CreatedAt, usersServices.CreatedBy, usersServices.ModifiedAt, usersServices.ModifiedBy).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := commands.NewUsersServicesCommand(db).InsertOne(context.TODO(), usersServices)
		assert.NoError(t, err)
	})
	t.Run("Error Begin", func(t *testing.T) {
		db, mock := NewMock()

		query := regexp.QuoteMeta("INSERT INTO users_services ( id, user_id, service_id, created_at, created_by, modified_at, modified_by ) VALUES ( ?, ?, ?, ?, ?, ?, ? );")

		mock.ExpectExec(query).WithArgs(usersServices.ID, nil, usersServices.ServiceID, usersServices.CreatedAt, usersServices.CreatedBy, usersServices.ModifiedAt, usersServices.ModifiedBy).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := commands.NewUsersServicesCommand(db).InsertOne(context.TODO(), usersServices)
		assert.Error(t, err)
	})
	t.Run("Error Exec", func(t *testing.T) {
		db, mock := NewMock()

		query := regexp.QuoteMeta("INSERT INTO users_services ( id, user_id, service_id, created_at, created_by, modified_at, modified_by ) VALUES ( ?, ?, ?, ?, ?, ?, ? );")

		mock.ExpectBegin()
		mock.ExpectExec(query).WithArgs(usersServices.ID, nil, usersServices.ServiceID, usersServices.CreatedAt, usersServices.CreatedBy, usersServices.ModifiedAt, usersServices.ModifiedBy).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := commands.NewUsersServicesCommand(db).InsertOne(context.TODO(), usersServices)
		assert.Error(t, err)
	})
	t.Run("Error Commit", func(t *testing.T) {
		db, mock := NewMock()

		query := regexp.QuoteMeta("INSERT INTO users_services ( id, user_id, service_id, created_at, created_by, modified_at, modified_by ) VALUES ( ?, ?, ?, ?, ?, ?, ? );")

		mock.ExpectBegin()
		mock.ExpectExec(query).WithArgs(usersServices.ID, usersServices.UserID, usersServices.ServiceID, usersServices.CreatedAt, usersServices.CreatedBy, usersServices.ModifiedAt, usersServices.ModifiedBy).WillReturnResult(sqlmock.NewResult(0, 0))

		err := commands.NewUsersServicesCommand(db).InsertOne(context.TODO(), usersServices)
		assert.Error(t, err)
	})
	t.Run("DB no connection", func(t *testing.T) {
		db, _ := NewMock()

		db.Close()
		err := commands.NewUsersServicesCommand(db).InsertOne(context.TODO(), usersServices)
		assert.Error(t, err)
		assert.EqualError(t, err, "error establishing a database connection")
	})
}
