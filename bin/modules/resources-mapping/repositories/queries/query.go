package queries

import (
	"context"
	models "engine/bin/modules/resources-mapping/models/domain"
)

type ResourcesMappingPostgre interface {
	FindByServiceId(ctx context.Context, serviceId string) (resourcesMappingList models.ResourcesMappingList, err error)
}
