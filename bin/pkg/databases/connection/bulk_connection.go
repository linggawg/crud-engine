package connection

import (
	dbsModels "engine/bin/modules/dbs/models/domain"
	"engine/bin/pkg/databases"
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
)

type BulkConnection struct {
	Sql map[string]*sqlx.DB
}

type BulkConnectionPkg struct {
	Connection BulkConnection
}

var pkgBulkConnection BulkConnectionPkg

func InitBulkConnectionPkg() *BulkConnectionPkg {
	if pkgBulkConnection.Connection.Sql == nil {
		pkgBulkConnection = BulkConnectionPkg{}
		pkgBulkConnection.Connection.Sql = make(map[string]*sqlx.DB)
	}
	return &pkgBulkConnection
}
func (h *BulkConnectionPkg) GetBulkConnectionSql(dbs dbsModels.Dbs) (val *sqlx.DB, err error) {
	val, ok := h.Connection.Sql[dbs.ID]
	if !ok {
		val, err = h.createConnection(dbs)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if _, ok := h.Connection.Sql[dbs.ID]; ok {
			log.Println(dbs.ID + " already exist")
			return nil, errors.New("cannot add connection to bulk manager connection")
		}
		h.Connection.Sql[dbs.ID] = val
	}
	return val, nil
}

func (h *BulkConnectionPkg) createConnection(dbs dbsModels.Dbs) (database *sqlx.DB, err error) {
	database, err = databases.InitDbs(databases.SQLXConfig{
		Host:     dbs.Host,
		Port:     uint16(dbs.Port),
		Name:     dbs.Name,
		Username: dbs.Username,
		Password: func() string {
			if dbs.Password != nil {
				return *dbs.Password
			} else {
				return ""
			}
		}(),
		Dialect: dbs.Dialect,
	})
	return database, err
}
