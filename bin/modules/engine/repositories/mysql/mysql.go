package mysql

import (
	"context"
	sqld "database/sql"
	"encoding/json"
	models "engine/bin/modules/engine/models/domain"
	"engine/bin/pkg/utils"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MySQL struct {
}

func (s *MySQL) FindData(ctx context.Context, conn interface{}, query string) ([]map[string]interface{}, error) {
	err := conn.(*sqlx.DB).PingContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error establishing a database connection")
	}

	rows, err := conn.(*sqlx.DB).QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				err = json.Unmarshal(b, &v)
				if err != nil {
					v = string(b)
				}
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}

	_, err = json.Marshal(tableData)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return tableData, nil
}

func (s *MySQL) CountData(ctx context.Context, conn interface{}, param string) (total int64, err error) {
	query := fmt.Sprintf(utils.QueryGetCount, param)
	err = conn.(*sqlx.DB).PingContext(ctx)
	if err != nil {
		log.Println(err)
		return 0, errors.New("error establishing a database connection")
	}

	err = conn.(*sqlx.DB).QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return total, nil
}

func (s *MySQL) FindPrimaryKey(ctx context.Context, conn interface{}, table string) (key *models.PrimaryKey, err error) {
	var (
		sql      string
		keyMysql models.KeyMySQL
	)

	err = conn.(*sqlx.DB).PingContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error establishing a database connection")
	}

	//sql = `show columns from ? where 'Key' = ? ;`
	sql = "show columns from " + table + " where `Key` = ?;"
	err = conn.(*sqlx.DB).GetContext(ctx, &keyMysql, sql, "PRI")
	if err != nil {
		if err != sqld.ErrNoRows {
			log.Println(err)
		}
		return nil, err
	}

	key = &models.PrimaryKey{
		Column: keyMysql.Field,
		Format: func() string {
			if strings.Contains(keyMysql.TypeData, "int") {
				return "int"
			} else {
				return "varchar"
			}
		}(),
	}
	return key, nil
}

func (s *MySQL) SelectInformationSchema(ctx context.Context, conn interface{}, table string) (informationSchema []models.InformationSchema, err error) {
	var (
		is   []models.InformationSchema
		sql  string
		args []interface{}
	)
	err = conn.(*sqlx.DB).PingContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error establishing a database connection")
	}

	sql = `SELECT column_name, is_nullable FROM information_schema.COLUMNS WHERE table_schema = DATABASE () AND table_name = ? ORDER BY table_name, ordinal_position;`
	args = append(args, table)
	err = conn.(*sqlx.DB).SelectContext(ctx, &is, sql, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return is, nil
}

func (s *MySQL) InsertOne(ctx context.Context, conn interface{}, query string, args []interface{}) (err error) {
	err = conn.(*sqlx.DB).PingContext(ctx)
	if err != nil {
		log.Println(err)
		return errors.New("error establishing a database connection")
	}

	_, err = conn.(*sqlx.DB).ExecContext(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *MySQL) UpdateOne(ctx context.Context, conn interface{}, query string, args []interface{}) (err error) {
	err = conn.(*sqlx.DB).PingContext(ctx)
	if err != nil {
		log.Println(err)
		return errors.New("error establishing a database connection")
	}

	_, err = conn.(*sqlx.DB).ExecContext(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *MySQL) DeleteOne(ctx context.Context, conn interface{}, query string, args []interface{}) (err error) {
	err = conn.(*sqlx.DB).PingContext(ctx)
	if err != nil {
		log.Println(err)
		return errors.New("error establishing a database connection")
	}

	_, err = conn.(*sqlx.DB).ExecContext(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
