package repositories

import (
	"context"
	models "engine/bin/modules/engine/models/domain"
)

type Repository interface {
	FindData(ctx context.Context, conn interface{}, query string) ([]map[string]interface{}, error)
	CountData(ctx context.Context, conn interface{}, param string) (total int64, err error)
	FindPrimaryKey(ctx context.Context, conn interface{}, table string) (key *models.PrimaryKey, err error)
	SelectInformationSchema(ctx context.Context, conn interface{}, table string) (informationSchema []models.InformationSchema, err error)

	InsertOne(ctx context.Context, conn interface{}, query string, args []interface{}) (err error)
	UpdateOne(ctx context.Context, conn interface{}, query string, args []interface{}) (err error)
	DeleteOne(ctx context.Context, conn interface{}, query string, args []interface{}) (err error)
}
