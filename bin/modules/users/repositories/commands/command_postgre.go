package commands

import (
	"context"
	models "engine/bin/modules/users/models/domain"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type UsersPostgreCommand struct {
	db *sqlx.DB
}

func NewUsersCommand(db *sqlx.DB) *UsersPostgreCommand {
	return &UsersPostgreCommand{db}
}

func (s *UsersPostgreCommand) InsertOne(ctx context.Context, users *models.Users) error {
	query := `
	INSERT INTO users
		(
		 	id,
			role_id,
			username,
			password,
			created_at,
			created_by,
		 	modified_at,
			modified_by
		) 
		VALUES 
		(
		 	:id,
			:role_id,
			:username,
			:password,
			:created_at,
			:created_by,
		 	:modified_at,
			:modified_by
		);
	`

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return errors.New("error establishing a database connection")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}
	defer tx.Rollback()

	_, err = s.db.NamedExecContext(ctx, query, users)
	if err != nil {
		log.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return errors.New("failed insert user")
	}

	return nil
}
