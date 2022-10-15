package commands

import (
	"context"
	models "engine/bin/modules/users-services/models/domain"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type UsersServicesPostgreCommand struct {
	db *sqlx.DB
}

func NewUsersServicesCommand(db *sqlx.DB) *UsersServicesPostgreCommand {
	return &UsersServicesPostgreCommand{db}
}

func (s *UsersServicesPostgreCommand) InsertOne(ctx context.Context, usersServices *models.UsersServices) error {
	query := `
	INSERT INTO users_services
		(
		 	id,
			user_id,
			service_id,
			created_at,
			created_by,
		 	modified_at,
			modified_by
		) 
		VALUES 
		(
		 	:id,
			:user_id,
			:service_id,
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

	_, err = s.db.NamedExecContext(ctx, query, usersServices)
	if err != nil {
		log.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return errors.New("failed insert users services")
	}

	return nil
}

func (s *UsersServicesPostgreCommand) DeleteByUsersIdAndServiceUrl(ctx context.Context, userId, serviceUrl string) error {
	query := `
	DELETE FROM users_services 
	WHERE user_id = $1 
	AND service_id IN (SELECT id FROM services WHERE service_url = $2);
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

	_, err = s.db.ExecContext(ctx, query, userId, serviceUrl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return errors.New("failed delete users services")
	}

	return nil
}
