package commands

import (
	"context"
	models "engine/bin/modules/users/models/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserspostgreCommand struct {
	db *sqlx.DB
}

func NewUsersCommand(db *sqlx.DB) *UserspostgreCommand {
	return &UserspostgreCommand{db}
}

func (s *UserspostgreCommand) InsertOne(ctx context.Context, users *models.Users) error {
	query := `
	INSERT INTO users
		(
		 	id,
			username,
			email,
			password,
			created_at,
			created_by,
		 	modified_at,
			modified_by
		) 
		VALUES 
		(
		 	:id,
			:username,
			:email,
			:password,
			:created_at,
			:created_by,
		 	:modified_at,
			:modified_by
		) RETURNING id ;
	`
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := s.db.NamedQueryContext(ctx, query, &models.Users{
		ID:         uuid.New().String(),
		Username:   users.Username,
		Email:      users.Email,
		Password:   users.Password,
		CreatedAt:  users.CreatedAt,
		CreatedBy:  users.CreatedBy,
		ModifiedAt: users.CreatedAt,
		ModifiedBy: &users.CreatedBy,
	})
	if err != nil {
		return err
	}
	if res.Next() {
		res.Scan(&users.ID)
	}
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}
