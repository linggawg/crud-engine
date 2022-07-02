package queries

import (
	"context"
	models "engine/bin/modules/engine/models/domain"
	"github.com/jmoiron/sqlx"
)

type EngineSQL interface {
	FindData(ctx context.Context, db *sqlx.DB, query string) ([]map[string]interface{}, error)
	CountData(ctx context.Context, db *sqlx.DB, param string) (total int64, err error)
	FindPrimaryKey(ctx context.Context, db *sqlx.DB, dialect, table string) (key *models.PrimaryKey, err error)
	SelectInformationSchema(ctx context.Context, db *sqlx.DB, dialect, table string) (informationSchema []models.InformationSchema, err error)
}
