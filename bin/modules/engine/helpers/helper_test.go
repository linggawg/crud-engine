package helpers_test

import (
	"engine/bin/modules/engine/helpers"
	resourcesMappingModels "engine/bin/modules/resources-mapping/models/domain"
	"engine/bin/pkg/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestSetQuery(t *testing.T) {
	t.Run("success-postgres", func(t *testing.T) {
		res := helpers.SetQuery(utils.DialectPostgres, "SELECT * FROM users WHERE id = ? AND name = ?")
		assert.Equal(t, res, "SELECT * FROM users WHERE id = $1 AND name = $2")
	})
	t.Run("success-mysql", func(t *testing.T) {
		res := helpers.SetQuery(utils.DialectMysql, "SELECT * FROM users WHERE id = ? AND name = ?")
		assert.Equal(t, res, "SELECT * FROM users WHERE id = ? AND name = ?")
	})
}

func TestConvertToSourceAlias(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var random = uuid.New().String()
		targetMap := make(map[string]interface{})
		targetMap["code_id"] = "a"
		targetMap["code_name"] = "name_a"

		originMap := make(map[string]interface{})
		originMap["code"] = "code_a"
		originMap["id"] = "a"
		originMap["name"] = "name_a"
		originMap["parent_id"] = "parent_code_a"

		var resourcesMappingList resourcesMappingModels.ResourcesMappingList
		resourcesMappingList = append(resourcesMappingList, resourcesMappingModels.ResourcesMapping{
			ID:           random,
			ServiceId:    random,
			SourceOrigin: "id",
			SourceAlias:  "code_id",
			CreatedBy:    &random,
			CreatedAt:    null.TimeFrom(time.Now()),
			ModifiedAt:   null.TimeFrom(time.Now()),
			ModifiedBy:   &random,
		})
		resourcesMappingList = append(resourcesMappingList, resourcesMappingModels.ResourcesMapping{
			ID:           random,
			ServiceId:    random,
			SourceOrigin: "name",
			SourceAlias:  "code_name",
			CreatedBy:    &random,
			CreatedAt:    null.TimeFrom(time.Now()),
			ModifiedAt:   null.TimeFrom(time.Now()),
			ModifiedBy:   &random,
		})
		res := helpers.ConvertToSourceAlias(originMap, resourcesMappingList)
		assert.Equal(t, res, targetMap)
	})
}

func TestConvertToSourceOrigin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var random = uuid.New().String()
		targetMap := make(map[string]interface{})
		targetMap["name"] = "name_c"
		targetMap["parent_id"] = "parent_code_c"

		aliasMap := make(map[string]interface{})
		aliasMap["code"] = "code_c"
		aliasMap["code_name"] = "name_c"
		aliasMap["parent_id"] = "parent_code_c"

		var resourcesMappingList resourcesMappingModels.ResourcesMappingList
		resourcesMappingList = append(resourcesMappingList, resourcesMappingModels.ResourcesMapping{
			ID:           random,
			ServiceId:    random,
			SourceOrigin: "name",
			SourceAlias:  "code_name",
			CreatedBy:    &random,
			CreatedAt:    null.TimeFrom(time.Now()),
			ModifiedAt:   null.TimeFrom(time.Now()),
			ModifiedBy:   &random,
		})
		resourcesMappingList = append(resourcesMappingList, resourcesMappingModels.ResourcesMapping{
			ID:           random,
			ServiceId:    random,
			SourceOrigin: "parent_id",
			SourceAlias:  "parent_id",
			CreatedBy:    &random,
			CreatedAt:    null.TimeFrom(time.Now()),
			ModifiedAt:   null.TimeFrom(time.Now()),
			ModifiedBy:   &random,
		})
		res := helpers.ConvertToSourceOrigin(aliasMap, resourcesMappingList)
		assert.Equal(t, res, targetMap)
	})
}
