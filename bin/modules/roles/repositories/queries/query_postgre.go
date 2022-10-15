package queries

import (
	"context"
	models "engine/bin/modules/roles/models/domain"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type RolesPostgreQuery struct {
	db *sqlx.DB
}

func NewRolesQuery(db *sqlx.DB) *RolesPostgreQuery {
	return &RolesPostgreQuery{db}
}

func (r *RolesPostgreQuery) FindOneByID(ctx context.Context, id string) (roles *models.Roles, err error) {
	var role models.Roles
	query := `
	SELECT
		id,
		name,
		created_at,
		created_by,
		modified_at,
		modified_by
	FROM
		roles
	WHERE id = $1
	`

	err = r.db.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error establishing a database connection")
	}

	err = r.db.GetContext(ctx, &role, query, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &role, nil
}
