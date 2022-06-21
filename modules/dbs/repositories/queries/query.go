package queries

import (
	"context"
	models "crud-engine/modules/dbs/models/domain"
)

type DbsPostgre interface {
	GetByID(ctx context.Context, id string) (dbs *models.Dbs, err error)
}
