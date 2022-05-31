package dbs

import (
	"context"
	"crud-engine/modules/models"

	"github.com/jmoiron/sqlx"
)

type HttpSqlx struct {
	db *sqlx.DB
}

func (s *HttpSqlx) GetByID(ctx context.Context, id string) (dbs *models.Dbs, err error) {
	var db models.Dbs
	query := `
	SELECT
		id,
		app_id,
		name,
		host,
		username,
		password,
		dialect
	FROM
		dbs
	WHERE id = $1
		`
	err = s.db.GetContext(ctx, &db, query, id)
	if err != nil {
		return nil, err
	}
	return &db, nil
}