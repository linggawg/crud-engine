package queries

import (
	"context"
	models "engine/bin/modules/resources-mapping/models/domain"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type ResourceMappingPostgreQuery struct {
	db *sqlx.DB
}

func NewResourcesMappingQuery(db *sqlx.DB) *ResourceMappingPostgreQuery {
	return &ResourceMappingPostgreQuery{db}
}

func (s *ResourceMappingPostgreQuery) FindByServiceId(ctx context.Context, serviceId string) (resourcesMappingList models.ResourcesMappingList, err error) {
	var (
		sql  string
		args []interface{}
	)
	sql = `SELECT id, service_id, source_origin, source_alias FROM resources_mapping WHERE service_id = $1; `
	args = append(args, serviceId)

	err = s.db.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error establishing a database connection")
	}

	err = s.db.SelectContext(ctx, &resourcesMappingList, sql, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return resourcesMappingList, nil
}
