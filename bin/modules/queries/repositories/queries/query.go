package queries

import (
	"context"
	models "engine/bin/modules/queries/models/domain"
)

type QueriesPostgre interface {
	FindOneByKey(ctx context.Context, key string) (queries *models.Queries, err error)
}
