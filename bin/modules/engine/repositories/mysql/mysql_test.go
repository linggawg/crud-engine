package mysql_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"engine/bin/modules/engine/repositories"
	"engine/bin/modules/engine/repositories/mysql"
	"engine/bin/pkg/utils"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestMySQL(t *testing.T) {
	suite.Run(t, new(MySQLTest))
}

type MySQLTest struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository repositories.Repository
	db         *sqlx.DB
}

func (s *MySQLTest) SetupTest() {
	var (
		db  *sql.DB
		err error
	)
	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.repository = &mysql.MySQL{}
	s.db = sqlx.NewDb(db, "sqlmock")
}

func (s *MySQLTest) TestFindData() {
	query := "SELECT id, name, data FROM sample"
	s.Run("Success", func() {
		rows := sqlmock.NewRows([]string{"id", "name", "data"})
		rows.AddRow("53b4aac9-fae2-4f54-b7c1-f88d8caef32f", "test 1", json.RawMessage(`{ "success": true, "data": [] }`))
		rows.AddRow("0dcb316f-6919-445d-a036-bc621875acc9", "test 2", json.RawMessage(`{ "success": true, "data": [] }`))
		s.mock.ExpectQuery(query).WillReturnRows(rows)
		data, err := s.repository.FindData(context.TODO(), s.db, query)
		s.Equal(reflect.TypeOf(make(map[string]interface{})), reflect.TypeOf(data[0]["data"]))
		s.NoError(err)
		s.Len(data, 2)
	})
	s.Run("JSON_Error", func() {
		rows := sqlmock.NewRows([]string{"id", "name", "data"})
		rows.AddRow("53b4aac9-fae2-4f54-b7c1-f88d8caef32f", "test 2", json.RawMessage(`{ "success": false, "data": ["not json format"] `))
		s.mock.ExpectQuery(query).WillReturnRows(rows)
		data, err := s.repository.FindData(context.TODO(), s.db, query)
		s.NoError(err)
		s.Equal(reflect.TypeOf(""), reflect.TypeOf(data[0]["data"]))
		s.Len(data, 1)
	})
	s.Run("Error", func() {
		s.mock.ExpectQuery(query).WillReturnError(sql.ErrConnDone)
		data, err := s.repository.FindData(context.TODO(), s.db, query)
		s.Empty(data)
		s.Error(err)
	})
	s.Run("DB no connection", func() {
		s.db.Close()
		s.mock.ExpectQuery(query).WillReturnError(driver.ErrBadConn)
		data, err := s.repository.FindData(context.TODO(), s.db, query)
		s.Empty(data)
		s.Error(err)
		s.EqualError(err, "error establishing a database connection")
	})
}

func (s *MySQLTest) TestCountData() {
	var (
		count int64 = 10
		param       = "SELECT * FROM users_services"
		query       = regexp.QuoteMeta(fmt.Sprintf(utils.QueryGetCount, param))
	)
	s.Run("Success", func() {
		s.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(count))
		total, err := s.repository.CountData(context.TODO(), s.db, param)
		s.NoError(err)
		s.Equal(total, count)
	})
	s.Run("Error", func() {
		s.mock.ExpectQuery(query).WillReturnError(sql.ErrConnDone)
		total, err := s.repository.CountData(context.TODO(), s.db, param)
		s.NotEqual(total, count)
		s.Error(err, sql.ErrConnDone)
	})
	s.Run("DB no connection", func() {
		s.db.Close()
		s.mock.ExpectQuery(query).WillReturnError(driver.ErrBadConn)
		data, err := s.repository.CountData(context.TODO(), s.db, query)
		s.Empty(data)
		s.Error(err)
		s.EqualError(err, "error establishing a database connection")
	})
}

func (s *MySQLTest) TestFindPrimaryKey() {
	rows := sqlmock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"})
	query := regexp.QuoteMeta("show columns from users where `Key` = ?;")
	s.Run("Success Primary Key Varchar", func() {
		rows.AddRow("id", "varchar(255)", "NO", "PRI", nil, nil)
		s.mock.ExpectQuery(query).WithArgs("PRI").WillReturnRows(rows)
		data, err := s.repository.FindPrimaryKey(context.TODO(), s.db, "users")
		s.NoError(err)
		s.Equal(data.Format, "varchar")
		s.NotNil(data)
	})
	s.Run("Success Primary Key Int", func() {
		rows.AddRow("id", "int", "NO", "PRI", nil, nil)
		s.mock.ExpectQuery(query).WithArgs("PRI").WillReturnRows(rows)
		data, err := s.repository.FindPrimaryKey(context.TODO(), s.db, "users")
		s.NoError(err)
		s.Equal(data.Format, "int")
		s.NotNil(data)
	})
	s.Run("ErrNoRows", func() {
		s.mock.ExpectQuery(query).WithArgs("PRI").WillReturnRows(rows)
		data, err := s.repository.FindPrimaryKey(context.TODO(), s.db, "users")
		s.Empty(data)
		s.Error(err, sql.ErrNoRows)
	})
	s.Run("DB no connection", func() {
		s.db.Close()
		s.mock.ExpectQuery(query).WillReturnError(driver.ErrBadConn)
		data, err := s.repository.FindPrimaryKey(context.TODO(), s.db, query)
		s.Empty(data)
		s.Error(err)
		s.EqualError(err, "error establishing a database connection")
	})
}

func (s *MySQLTest) TestSelectInformationSchema() {
	rows := sqlmock.NewRows([]string{"column_name", "is_nullable"})
	rows.AddRow("id", "NO")
	rows.AddRow("name", "YES")
	query := regexp.QuoteMeta("SELECT column_name, is_nullable FROM information_schema.COLUMNS WHERE table_schema = DATABASE () AND table_name = ? ORDER BY table_name, ordinal_position;")
	s.Run("Success", func() {
		s.mock.ExpectQuery(query).WithArgs("users").WillReturnRows(rows)
		data, err := s.repository.SelectInformationSchema(context.TODO(), s.db, "users")
		s.NoError(err)
		s.Len(data, 2)
	})
	s.Run("RowsWillBeClosed", func() {
		s.mock.ExpectQuery(query).WithArgs("users").RowsWillBeClosed()
		data, err := s.repository.SelectInformationSchema(context.TODO(), s.db, "users")
		s.Error(err)
		s.Empty(data)
	})
	s.Run("DB no connection", func() {
		s.db.Close()
		s.mock.ExpectQuery(query).WillReturnError(driver.ErrBadConn)
		data, err := s.repository.SelectInformationSchema(context.TODO(), s.db, query)
		s.Empty(data)
		s.Error(err)
		s.EqualError(err, "error establishing a database connection")
	})
}

func (s *MySQLTest) TestInsertOne() {
	query := regexp.QuoteMeta("INSERT INTO sample ( id, created_at ) VALUES ( ?, ?);")
	var args []interface{}
	args = append(args, "0dcb316f-6919-445d-a036-bc621875acc9")
	args = append(args, "2022-07-09T15:50:26+07:00")
	s.Run("Success", func() {
		s.mock.ExpectExec(query).WithArgs("0dcb316f-6919-445d-a036-bc621875acc9", "2022-07-09T15:50:26+07:00").WillReturnResult(sqlmock.NewResult(0, 1))
		err := s.repository.InsertOne(context.TODO(), s.db, "INSERT INTO sample ( id, created_at ) VALUES ( ?, ?);", args)
		s.NoError(err)
	})
	s.Run("Error", func() {
		s.mock.ExpectExec(query).WithArgs("0dcb316f-6919-445d-a036-bc621875acc9", "2022-07-09T15:50:26+07:00").WillReturnResult(sqlmock.NewResult(0, 0))
		err := s.repository.InsertOne(context.TODO(), s.db, query, args)
		s.Error(err)
	})
	s.Run("DB no connection", func() {
		s.db.Close()
		s.mock.ExpectQuery(query).WillReturnError(driver.ErrBadConn)
		err := s.repository.InsertOne(context.TODO(), s.db, query, args)
		s.Error(err)
		s.EqualError(err, "error establishing a database connection")
	})
}

func (s *MySQLTest) TestUpdateOne() {
	query := regexp.QuoteMeta("UPDATE sample SET modified_at = ? WHERE id = ?;")
	var args []interface{}
	args = append(args, "0dcb316f-6919-445d-a036-bc621875acc9")
	args = append(args, "2022-07-09T15:50:26+07:00")
	s.Run("Success", func() {
		s.mock.ExpectExec(query).WithArgs("0dcb316f-6919-445d-a036-bc621875acc9", "2022-07-09T15:50:26+07:00").WillReturnResult(sqlmock.NewResult(0, 1))
		err := s.repository.UpdateOne(context.TODO(), s.db, "UPDATE sample SET modified_at = ? WHERE id = ?;", args)
		s.NoError(err)
	})
	s.Run("Error", func() {
		s.mock.ExpectExec(query).WithArgs("0dcb316f-6919-445d-a036-bc621875acc9", "2022-07-09T15:50:26+07:00").WillReturnResult(sqlmock.NewResult(0, 0))
		err := s.repository.UpdateOne(context.TODO(), s.db, query, args)
		s.Error(err)
	})
	s.Run("DB no connection", func() {
		s.db.Close()
		s.mock.ExpectQuery(query).WillReturnError(driver.ErrBadConn)
		err := s.repository.UpdateOne(context.TODO(), s.db, query, args)
		s.Error(err)
		s.EqualError(err, "error establishing a database connection")
	})
}

func (s *MySQLTest) TestDeleteOne() {
	query := regexp.QuoteMeta("DELETE FROM sample WHERE id = ?;")
	var args []interface{}
	args = append(args, "0dcb316f-6919-445d-a036-bc621875acc9")
	s.Run("Success", func() {
		s.mock.ExpectExec(query).WithArgs("0dcb316f-6919-445d-a036-bc621875acc9").WillReturnResult(sqlmock.NewResult(0, 1))
		err := s.repository.DeleteOne(context.TODO(), s.db, "DELETE FROM sample WHERE id = ?;", args)
		s.NoError(err)
	})
	s.Run("Error", func() {
		s.mock.ExpectExec(query).WithArgs("0dcb316f-6919-445d-a036-bc621875acc9").WillReturnResult(sqlmock.NewResult(0, 0))
		err := s.repository.DeleteOne(context.TODO(), s.db, query, args)
		s.Error(err)
	})
	s.Run("DB no connection", func() {
		s.db.Close()
		s.mock.ExpectQuery(query).WillReturnError(driver.ErrBadConn)
		err := s.repository.DeleteOne(context.TODO(), s.db, query, args)
		s.Error(err)
		s.EqualError(err, "error establishing a database connection")
	})
}
