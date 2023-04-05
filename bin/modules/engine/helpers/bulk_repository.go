package helpers

import (
	"engine/bin/modules/engine/repositories"
	"engine/bin/modules/engine/repositories/mysql"
	"engine/bin/modules/engine/repositories/postgresql"
	"engine/bin/pkg/utils"
)

type BulkRepository struct {
	repo map[string]repositories.Repository
}

func InitBulkRepository() *BulkRepository {
	var repo = make(map[string]repositories.Repository)
	repo[utils.DialectPostgres] = &postgresql.PostgreSQL{}
	repo[utils.DialectMysql] = &mysql.MySQL{}
	return &BulkRepository{repo}
}

func (h *BulkRepository) GetBulkRepository(id string) repositories.Repository {
	val, ok := h.repo[id]
	if !ok {
		panic(id + " not found")
	}
	return val
}
