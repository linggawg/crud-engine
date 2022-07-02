package helpers

import (
	dbsModels "engine/bin/modules/dbs/models/domain"
	conn "engine/bin/pkg/databases"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreateConnection(dbs dbsModels.Dbs) (database *sqlx.DB, err error) {
	database, err = conn.InitDbs(conn.SQLXConfig{
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
func SetQuery(dialect, query string) (newQuery string) {
	switch dialect {
	case "postgres":
		count := strings.Count(query, "?")
		for i := 1; i <= count; i++ {
			query = strings.Replace(query, "?", fmt.Sprintf("$%d", i), 1)
		}
		return query
	default:
		return query
	}
}
