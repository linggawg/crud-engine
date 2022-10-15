package commands

import (
	"context"
	models "engine/bin/modules/services/models/domain"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type ServicesPostgreCommand struct {
	db *sqlx.DB
}

func NewServicesCommand(db *sqlx.DB) *ServicesPostgreCommand {
	return &ServicesPostgreCommand{db}
}

func (s *ServicesPostgreCommand) InsertOne(ctx context.Context, services *models.Services) error {
	query := `
	INSERT INTO services
		(
			id,
			db_id,
			query_id,
			method,
			service_url,
			created_at,
			created_by,
		 	modified_at,
			modified_by
		) 
		VALUES 
		(
		 	:id,
			:db_id,
			:query_id,
			:method,
			:service_url,
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

	_, err = s.db.NamedExecContext(ctx, query, services)
	if err != nil {
		log.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return errors.New("failed insert services")
	}

	return nil
}

func (s *ServicesPostgreCommand) DeleteByServiceUrl(ctx context.Context, serviceUrl string) error {
	query := `DELETE FROM services WHERE service_url = $1;`

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

	_, err = s.db.ExecContext(ctx, query, serviceUrl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return errors.New("failed delete services")
	}

	return nil
}
