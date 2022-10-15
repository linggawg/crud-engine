package queries

import (
	"context"
	models "engine/bin/modules/dbs/models/domain"
)

type DbsPostgre interface {
	FindOneByID(ctx context.Context, id string) (dbs *models.Dbs, err error)
}
