package queries

import (
	"context"
	models "engine/bin/modules/queries/models/domain"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type QueriesPostgreQuery struct {
	db *sqlx.DB
}

func NewQueriesQuery(db *sqlx.DB) *QueriesPostgreQuery {
	return &QueriesPostgreQuery{db}
}

func (q *QueriesPostgreQuery) FindOneByKey(ctx context.Context, key string) (queries *models.Queries, err error) {
	var queriesModel models.Queries
	query := `
	SELECT
		id,
		key,
		query_definition,
		created_at,
		created_by,
		modified_at,
		modified_by
	FROM
		queries
	WHERE key = $1
	`

	err = q.db.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error establishing a database connection")
	}

	err = q.db.GetContext(ctx, &queriesModel, query, key)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &queriesModel, nil
}
