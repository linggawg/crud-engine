package helpers

import (
	"engine/bin/modules/engine/repositories/mysql"
	"engine/bin/modules/engine/repositories/postgresql"
	"engine/bin/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetBulkRepositoryy(t *testing.T) {
	bulkRepository := InitBulkRepository()
	t.Run("success-repo-postgres", func(t *testing.T) {
		res := bulkRepository.GetBulkRepository(utils.DialectPostgres)
		assert.Equal(t, res, &postgresql.PostgreSQL{})
	})
	t.Run("success-repo-mysql", func(t *testing.T) {
		res := bulkRepository.GetBulkRepository(utils.DialectMysql)
		assert.Equal(t, res, &mysql.MySQL{})
	})
	t.Run("error-repo-not-found", func(t *testing.T) {
		t.Helper()
		defer func() { _ = recover() }()
		bulkRepository.GetBulkRepository("postgres_1")
		t.Errorf("should have panicked")
	})
}
