package databases

import (
	"errors"
	"github.com/jmoiron/sqlx"
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

func (h *BulkConnectionPkg) AddBulkConnectionSql(id string, value *sqlx.DB) error {
	if _, ok := h.Connection.Sql[id]; ok {
		return errors.New(id + " already exist")
	}
	h.Connection.Sql[id] = value
	return nil
}

func (h *BulkConnectionPkg) GetBulkConnectionSql(id string) (*sqlx.DB, error) {
	val, ok := h.Connection.Sql[id]
	if !ok {
		return nil, errors.New(id + " not found")
	}
	return val, nil
}
