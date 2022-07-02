package commands

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type EngineSQL interface {
	InsertOne(ctx context.Context, db *sqlx.DB, query string, args []interface{}) (err error)
	UpdateOne(ctx context.Context, db *sqlx.DB, query string, args []interface{}) (err error)
	DeleteOne(ctx context.Context, db *sqlx.DB, query string, args []interface{}) (err error)
}
