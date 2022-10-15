package commands_test

import (
	"context"
	models "engine/bin/modules/users/models/domain"
	"engine/bin/modules/users/repositories/commands"
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
var users = &models.Users{
	ID:         random,
	RoleID:     random,
	Username:   random,
	Password:   random,
	CreatedAt:  null.TimeFrom(time.Now()),
	CreatedBy:  &random,
	ModifiedAt: null.TimeFrom(time.Now()),
	ModifiedBy: &random,
}

func TestInsertOne(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock := NewMock()

		query := regexp.QuoteMeta("INSERT INTO users ( id, role_id, username, password, created_at, created_by, modified_at, modified_by ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ? );")

		mock.ExpectBegin()
		mock.ExpectExec(query).WithArgs(users.ID, users.RoleID, users.Username, users.Password, users.CreatedAt, users.CreatedBy, users.ModifiedAt, users.ModifiedBy).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := commands.NewUsersCommand(db).InsertOne(context.TODO(), users)
		assert.NoError(t, err)
	})
	t.Run("Error Begin", func(t *testing.T) {
		db, mock := NewMock()

		query := regexp.QuoteMeta("INSERT INTO users ( id, role_id, username, password, created_at, created_by, modified_at, modified_by ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ? );")

		mock.ExpectExec(query).WithArgs(users.ID, users.RoleID, users.Username, users.Password, users.CreatedAt, users.CreatedBy, users.ModifiedAt, users.ModifiedBy).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := commands.NewUsersCommand(db).InsertOne(context.TODO(), users)
		assert.Error(t, err)
	})
	t.Run("Error Exec", func(t *testing.T) {
		db, mock := NewMock()

		query := regexp.QuoteMeta("INSERT INTO users ( id, role_id, username, password, created_at, created_by, modified_at, modified_by ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ? );")

		mock.ExpectBegin()
		mock.ExpectExec(query).WithArgs(users.ID, nil, users.Username, users.Password, users.CreatedAt, users.CreatedBy, users.ModifiedAt, users.ModifiedBy).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := commands.NewUsersCommand(db).InsertOne(context.TODO(), users)
		assert.Error(t, err)
	})
	t.Run("Error Commit", func(t *testing.T) {
		db, mock := NewMock()

		query := regexp.QuoteMeta("INSERT INTO users ( id, role_id, username, password, created_at, created_by, modified_at, modified_by ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ? );")

		mock.ExpectBegin()
		mock.ExpectExec(query).WithArgs(users.ID, users.RoleID, users.Username, users.Password, users.CreatedAt, users.CreatedBy, users.ModifiedAt, users.ModifiedBy).WillReturnResult(sqlmock.NewResult(0, 0))

		err := commands.NewUsersCommand(db).InsertOne(context.TODO(), users)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed insert user")
	})
	t.Run("DB no connection", func(t *testing.T) {
		db, _ := NewMock()

		db.Close()
		err := commands.NewUsersCommand(db).InsertOne(context.TODO(), users)
		assert.Error(t, err)
		assert.EqualError(t, err, "error establishing a database connection")
	})
}
