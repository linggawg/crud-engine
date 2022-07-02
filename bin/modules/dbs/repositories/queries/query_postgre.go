package queries

import (
	"context"
	models "engine/bin/modules/dbs/models/domain"

	"github.com/jmoiron/sqlx"
)

type DbsPostgreQuery struct {
	db *sqlx.DB
}

func NewDbsQuery(db *sqlx.DB) *DbsPostgreQuery {
	return &DbsPostgreQuery{db}
}

func (s *DbsPostgreQuery) GetByID(ctx context.Context, id string) (dbs *models.Dbs, err error) {
	var db models.Dbs
	query := `
			SELECT
				id,
				app_id,
				name,
				host,
				port,
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
