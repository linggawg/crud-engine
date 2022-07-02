package queries

import (
	"context"
	models "engine/bin/modules/users/models/domain"

	"github.com/jmoiron/sqlx"
)

type UsersPostgreQuery struct {
	db *sqlx.DB
}

func NewUsersQuery(db *sqlx.DB) *UsersPostgreQuery {
	return &UsersPostgreQuery{db}
}

func (s *UsersPostgreQuery) GetByID(ctx context.Context, id string) (users *models.Users, err error) {
	var u models.Users
	query := `
	SELECT
		id,
		username,
		email,
		password,
		created_at,
		created_by,
		modified_at,
		modified_by
	FROM
		users
	WHERE id = $1
		`
	err = s.db.GetContext(ctx, &u, query, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UsersPostgreQuery) GetByEmail(ctx context.Context, email string) (users *models.Users, err error) {
	var u models.Users
	query := `
	SELECT
		id,
		username,
		email,
		password,
		created_at,
		created_by,
		modified_at,
		modified_by
	FROM
		users
	WHERE email = $1
		`
	err = s.db.GetContext(ctx, &u, query, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
