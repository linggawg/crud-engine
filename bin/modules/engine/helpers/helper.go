package helpers

import (
	resourcesMappingModels "engine/bin/modules/resources-mapping/models/domain"
	"engine/bin/pkg/utils"
	"fmt"
	"strings"
)

func SetQuery(dialect, query string) (newQuery string) {
	switch dialect {
	case utils.DialectPostgres:
		count := strings.Count(query, "?")
		for i := 1; i <= count; i++ {
			query = strings.Replace(query, "?", fmt.Sprintf("$%d", i), 1)
		}
		return strings.ReplaceAll(query, "`", "\"")
	default:
		return query
	}
}

func ConvertToSourceAlias(originMap map[string]interface{}, resourcesMapping resourcesMappingModels.ResourcesMappingList) map[string]interface{} {
	targetMap := map[string]interface{}{}
	for key := range originMap {
		for _, rm := range resourcesMapping {
			if rm.SourceOrigin == key {
				targetMap[rm.SourceAlias] = originMap[key]
				continue
			}
		}
		continue
	}
	return targetMap
}

func ConvertToSourceOrigin(aliasMap map[string]interface{}, resourcesMapping resourcesMappingModels.ResourcesMappingList) map[string]interface{} {
	targetMap := map[string]interface{}{}
	for key := range aliasMap {
		for _, rm := range resourcesMapping {
			if rm.SourceAlias == key {
				targetMap[rm.SourceOrigin] = aliasMap[key]
				continue
			}
		}
		continue
	}
	return targetMap
}
