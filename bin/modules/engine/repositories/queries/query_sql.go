package queries

import (
	"context"
	"encoding/json"
	models "engine/bin/modules/engine/models/domain"
	"errors"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type EngineQuery struct {
}

func NewEngineQuery() *EngineQuery {
	return &EngineQuery{}
}
func (s *EngineQuery) FindData(ctx context.Context, db *sqlx.DB, query string) ([]map[string]interface{}, error) {
	rows, err := db.QueryContext(ctx, query)
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
				v = string(b)
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

func (s *EngineQuery) CountData(ctx context.Context, db *sqlx.DB, param string) (total int64, err error) {
	err = db.QueryRow("SELECT COUNT('total') FROM (" + param + ") as total").Scan(&total)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return total, err
}

func (s *EngineQuery) FindPrimaryKey(ctx context.Context, db *sqlx.DB, dialect, table string) (key *models.PrimaryKey, err error) {
	var sql string

	if dialect == "mysql" {
		var keyMysql models.KeyMySQL
		sql = `show columns from ? where 'Key' = ? ;`
		err = db.GetContext(ctx, &keyMysql, sql, table, "PRI")
		if err != nil {
			log.Println(err)
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
		return key, err
	} else if dialect == "postgres" {
		var keyPostgres models.KeyPostgres
		sql = `SELECT pg_attribute.attname, format_type(pg_attribute.atttypid, pg_attribute.atttypmod) 
		FROM pg_index, pg_class, pg_attribute, pg_namespace 
		WHERE pg_class.oid = $1::regclass 
		AND indrelid = pg_class.oid 
		AND pg_class.relnamespace = pg_namespace.oid 
		AND pg_attribute.attrelid = pg_class.oid 
		AND pg_attribute.attnum = any(pg_index.indkey) AND indisprimary`
		err = db.GetContext(ctx, &keyPostgres, sql, table)
		if err != nil {
			log.Println(err)
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
		return key, err
	}
	return nil, errors.New("does not support this dialect")

}

func (s *EngineQuery) SelectInformationSchema(ctx context.Context, db *sqlx.DB, dialect, table string) (informationSchema []models.InformationSchema, err error) {
	var (
		is   []models.InformationSchema
		sql  string
		args []interface{}
	)
	if dialect == "mysql" {
		sql = `SELECT column_name, is_nullable FROM information_schema.COLUMNS WHERE table_schema = DATABASE () AND table_name = ? ORDER BY table_name, ordinal_position;`
	} else if dialect == "postgres" {
		sql = `SELECT column_name, is_nullable FROM information_schema.columns WHERE table_schema = current_schema() AND table_name = $1 order by table_name,ordinal_position; `
	} else {
		return nil, errors.New("does not support this dialect")
	}
	args = append(args, table)
	err = db.SelectContext(ctx, &is, sql, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return is, nil
}
