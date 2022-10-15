package connection

import (
	dbsModels "engine/bin/modules/dbs/models/domain"
	"github.com/jmoiron/sqlx"
)

type Connection interface {
	GetBulkConnectionSql(dbs dbsModels.Dbs) (val *sqlx.DB, err error)
}
