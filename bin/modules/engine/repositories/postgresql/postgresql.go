package postgresql

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

type PostgreSQL struct {
}

func (s *PostgreSQL) FindData(ctx context.Context, conn interface{}, query string) ([]map[string]interface{}, error) {
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

func (s *PostgreSQL) CountData(ctx context.Context, conn interface{}, param string) (total int64, err error) {
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

func (s *PostgreSQL) FindPrimaryKey(ctx context.Context, conn interface{}, table string) (key *models.PrimaryKey, err error) {
	var (
		sql         string
		keyPostgres models.KeyPostgres
	)
	err = conn.(*sqlx.DB).PingContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error establishing a database connection")
	}

	sql = `SELECT pg_attribute.attname, format_type(pg_attribute.atttypid, pg_attribute.atttypmod) 
		FROM pg_index, pg_class, pg_attribute, pg_namespace 
		WHERE pg_class.oid = $1::regclass 
		AND indrelid = pg_class.oid 
		AND pg_class.relnamespace = pg_namespace.oid 
		AND pg_attribute.attrelid = pg_class.oid 
		AND pg_attribute.attnum = any(pg_index.indkey) AND indisprimary`
	err = conn.(*sqlx.DB).GetContext(ctx, &keyPostgres, sql, table)
	if err != nil {
		if err != sqld.ErrNoRows {
			log.Println(err)
		}
		return nil, err
	}

	key = &models.PrimaryKey{
		Column: keyPostgres.AttName,
		Format: func() string {
			if strings.Contains(keyPostgres.FormatType, "int") {
				return "int"
			} else {
				return "varchar"
			}
		}(),
	}
	return key, nil
}

func (s *PostgreSQL) SelectInformationSchema(ctx context.Context, conn interface{}, table string) (informationSchema []models.InformationSchema, err error) {
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

	sql = `SELECT column_name, is_nullable FROM information_schema.columns WHERE table_schema = current_schema() AND table_name = $1 order by table_name,ordinal_position; `
	args = append(args, table)
	err = conn.(*sqlx.DB).SelectContext(ctx, &is, sql, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return is, nil
}

func (s *PostgreSQL) InsertOne(ctx context.Context, conn interface{}, query string, args []interface{}) (err error) {
	err = conn.(*sqlx.DB).PingContext(ctx)
	if err != nil {
		log.Println(err)
		return errors.New("error establishing a database connection")
	}

	_, err = conn.(*sqlx.DB).ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("asdasddadsa")
		log.Println(err.Error())
		return err
	}
	return nil
}

func (s *PostgreSQL) UpdateOne(ctx context.Context, conn interface{}, query string, args []interface{}) (err error) {
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

func (s *PostgreSQL) DeleteOne(ctx context.Context, conn interface{}, query string, args []interface{}) (err error) {
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
