package queries

import (
	"context"
	models "engine/bin/modules/users/models/domain"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type UsersPostgreQuery struct {
	db *sqlx.DB
}

func NewUsersQuery(db *sqlx.DB) *UsersPostgreQuery {
	return &UsersPostgreQuery{db}
}

func (s *UsersPostgreQuery) FindOneByID(ctx context.Context, id string) (users *models.Users, err error) {
	var user models.Users
	query := `
	SELECT
		id,
		role_id,
		username,
		password,
		created_at,
		created_by,
		modified_at,
		modified_by
	FROM
		users
	WHERE id = $1
		`

	err = s.db.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error establishing a database connection")
	}

	err = s.db.GetContext(ctx, &user, query, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}

func (s *UsersPostgreQuery) FindOneByUsername(ctx context.Context, username string) (users *models.Users, err error) {
	var user models.Users
	query := `
	SELECT
		id,
		role_id,
		username,
		password,
		created_at,
		created_by,
		modified_at,
		modified_by
	FROM
		users
	WHERE username = $1
		`

	err = s.db.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error establishing a database connection")
	}

	err = s.db.GetContext(ctx, &user, query, username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}
